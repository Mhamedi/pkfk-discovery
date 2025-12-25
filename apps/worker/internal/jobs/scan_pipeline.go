package jobs
import (
	"context"
	""
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/pkfk-discovery/worker/internal/domain"
	"github.com/pkfk-discovery/worker/internal/integrations/minio"
	"github.com/pkfk-discovery/worker/internal/services"
)
type ScanPipeline struct {
	db           *pgxpool.Pool
	minioClient  *minio.Client
	logger       *logrus.Logger
	sqlSafety    *services.SQLSafety
}
type ScanJobPayload struct {
	ScanID uuid.UUID `json:"scan_id"`
}
type ScanResults struct {
	Metadata      MetadataResults      `json:"metadata"`
	Profiling     ProfilingResults     `json:"profiling"`
	Candidates    []Candidate          `json:"candidates"`
	Evidence      EvidenceResults      `json:"evidence"`
	Scoring       ScoringResults       `json:"scoring"`
	Graph         GraphResults         `json:"graph"`
	Report        ReportResults        `json:"report"`
}
type MetadataResults struct {
	Tables    []TableMetadata    `json:"tables"`
	Columns   []ColumnMetadata    `json:"columns"`
	Indexes   []IndexMetadata     `json:"indexes"`
	Constraints []ConstraintMetadata `json:"constraints"`
}
type TableMetadata struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}
type ColumnMetadata struct {
	Schema      string `json:"schema"`
	Table       string `json:"table"`
	Name        string `json:"name"`
	DataType    string `json:"data_type"`
	IsNullable  bool   `json:"is_nullable"`
	IsPrimaryKey bool  `json:"is_primary_key"`
}
type IndexMetadata struct {
	Schema    string   `json:"schema"`
	Table     string   `json:"table"`
	Name      string   `json:"name"`
	Columns   []string `json:"columns"`
	IsUnique  bool     `json:"is_unique"`
}
type ConstraintMetadata struct {
	Schema      string   `json:"schema"`
	Table       string   `json:"table"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // PRIMARY_KEY, FOREIGN_KEY, UNIQUE, CHECK
	Columns     []string `json:"columns"`
	References  *ReferenceMetadata `json:"references,omitempty"`
}
type ReferenceMetadata struct {
	Schema string   `json:"schema"`
	Table  string   `json:"table"`
	Columns []string `json:"columns"`
}
type ProfilingResults struct {
	ColumnStats map[string]ColumnStats `json:"column_stats"`
}
type ColumnStats struct {
	DistinctCount int64   `json:"distinct_count"`
	NullCount     int64   `json:"null_count"`
	SampleValues  []string `json:"sample_values,omitempty"`
}
type Candidate struct {
	FromSchema   string   `json:"from_schema"`
	FromTable    string   `json:"from_table"`
	FromColumns  []string `json:"from_columns"`
	ToSchema     string   `json:"to_schema"`
	ToTable      string   `json:"to_table"`
	ToColumns    []string `json:"to_columns"`
	Confidence   float64  `json:"confidence"`
	Reason       string   `json:"reason"`
}
type EvidenceResults struct {
	Evidence map[string][]EvidenceItem `json:"evidence"`
}
type EvidenceItem struct {
	Type        string  `json:"type"` // NAMING, VALUE_OVERLAP, CARDINALITY, INDEX
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}
type ScoringResults struct {
	Scores map[string]float64 `json:"scores"` // candidate_id -> score
}
type GraphResults struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}
type GraphNode struct {
	ID       string `json:"id"`
	Type     string `json:"type"` // TABLE, COLUMN
	Schema   string `json:"schema"`
	Table    string `json:"table"`
	Column   string `json:"column,omitempty"`
}
type GraphEdge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Type   string  `json:"type"` // PK_FK, FK_FK
	Confidence float64 `json:"confidence"`
}
type ReportResults struct {
	Summary      ReportSummary `json:"summary"`
	Recommendations []string   `json:"recommendations"`
}
type ReportSummary struct {
	TotalTables      int `json:"total_tables"`
	TotalColumns     int `json:"total_columns"`
	PKFKRelationships int `json:"pkfk_relationships"`
	HighConfidence   int `json:"high_confidence"`
	MediumConfidence int `json:"medium_confidence"`
	LowConfidence    int `json:"low_confidence"`
}
func NewScanPipeline(db *pgxpool.Pool, minioClient *minio.Client, logger *logrus.Logger, sqlSafety *services.SQLSafety) *ScanPipeline {
	return &ScanPipeline{
		db:          db,
		minioClient: minioClient,
		logger:      logger,
		sqlSafety:   sqlSafety,
	}
}
func (p *ScanPipeline) Run(ctx context.Context, scanID uuid.UUID, adapter *domain.Adapter, connection *domain.Connection, policy domain.ScanPolicy) (*ScanResults, error) {
	results := &ScanResults{}
	// Stage 1: Metadata Extraction
	p.logger.WithField("scan_id", scanID).Info("Starting metadata extraction")
	metadata, err := p.extractMetadata(ctx, adapter, connection)
	if err != nil {
		return nil, fmt.Errorf("metadata extraction failed: %w", err)
	}
	results.Metadata = *metadata
	// Stage 2: Profiling
	p.logger.WithField("scan_id", scanID).Info("Starting profiling")
	profiling, err := p.profile(ctx, adapter, connection, metadata, policy)
	if err != nil {
		return nil, fmt.Errorf("profiling failed: %w", err)
	}
	results.Profiling = *profiling
	// Stage 3: Candidate Generation
	p.logger.WithField("scan_id", scanID).Info("Starting candidate generation")
	candidates := p.generateCandidates(metadata, profiling)
	results.Candidates = candidates
	// Stage 4: Evidence Collection
	p.logger.WithField("scan_id", scanID).Info("Starting evidence collection")
	evidence, err := p.collectEvidence(ctx, adapter, connection, candidates, policy)
	if err != nil {
		return nil, fmt.Errorf("evidence collection failed: %w", err)
	}
	results.Evidence = *evidence
	// Stage 5: Scoring
	p.logger.WithField("scan_id", scanID).Info("Starting scoring")
	scoring := p.score(candidates, evidence)
	results.Scoring = *scoring
	// Stage 6: Graph Reconciliation
	p.logger.WithField("scan_id", scanID).Info("Starting graph reconciliation")
	graph := p.reconcileGraph(metadata, candidates, scoring)
	results.Graph = *graph
	// Stage 7: Report Generation
	p.logger.WithField("scan_id", scanID).Info("Generating report")
	report := p.generateReport(metadata, graph, scoring)
	results.Report = *report
	return results, nil
}
func (p *ScanPipeline) extractMetadata(ctx context.Context, adapter *domain.Adapter, connection *domain.Connection) (*MetadataResults, error) {
	// TODO: Load adapter bundle from MinIO and execute metadata SQL templates
	// For now, return placeholder structure
	return &MetadataResults{
		Tables:      []TableMetadata{},
		Columns:     []ColumnMetadata{},
		Indexes:     []IndexMetadata{},
		Constraints: []ConstraintMetadata{},
	}, nil
}
func (p *ScanPipeline) profile(ctx context.Context, adapter *domain.Adapter, connection *domain.Connection, metadata *MetadataResults, policy domain.ScanPolicy) (*ProfilingResults, error) {
	// TODO: Execute profiling SQL templates with sampling
	// For now, return placeholder structure
	return &ProfilingResults{
		ColumnStats: make(map[string]ColumnStats),
	}, nil
}
func (p *ScanPipeline) generateCandidates(metadata *MetadataResults, profiling *ProfilingResults) []Candidate {
	candidates := []Candidate{}
	// Generate PK candidates from primary key constraints
	for _, constraint := range metadata.Constraints {
		if constraint.Type == "PRIMARY_KEY" {
			// Find potential FK relationships
			for _, otherConstraint := range metadata.Constraints {
				if otherConstraint.Type == "FOREIGN_KEY" {
					// Already has FK constraint, skip
					continue
				}
				// Check if columns match naming patterns or data types
				// This is a simplified heuristic
				candidate := Candidate{
					FromSchema:  constraint.Schema,
					FromTable:   constraint.Table,
					FromColumns: constraint.Columns,
					ToSchema:    otherConstraint.Schema,
					ToTable:     otherConstraint.Table,
					ToColumns:   otherConstraint.Columns,
					Confidence:  0.3, // Low initial confidence
					Reason:      "Naming pattern match",
				}
				candidates = append(candidates, candidate)
			}
		}
	}
	return candidates
}
func (p *ScanPipeline) collectEvidence(ctx context.Context, adapter *domain.Adapter, connection *domain.Connection, candidates []Candidate, policy domain.ScanPolicy) (*EvidenceResults, error) {
	evidence := &EvidenceResults{
		Evidence: make(map[string][]EvidenceItem),
	}
	// TODO: Execute evidence collection SQL templates
	// For now, return placeholder structure
	for _, candidate := range candidates {
		candidateID := fmt.Sprintf("%s.%s.%s->%s.%s.%s",
			candidate.FromSchema, candidate.FromTable, candidate.FromColumns[0],
			candidate.ToSchema, candidate.ToTable, candidate.ToColumns[0])
		evidence.Evidence[candidateID] = []EvidenceItem{
			{
				Type:        "NAMING",
				Confidence:  0.5,
				Description: "Column naming pattern suggests relationship",
			},
		}
	}
	return evidence, nil
}
func (p *ScanPipeline) score(candidates []Candidate, evidence *EvidenceResults) *ScoringResults {
	scoring := &ScoringResults{
		Scores: make(map[string]float64),
	}
	for _, candidate := range candidates {
		candidateID := fmt.Sprintf("%s.%s.%s->%s.%s.%s",
			candidate.FromSchema, candidate.FromTable, candidate.FromColumns[0],
			candidate.ToSchema, candidate.ToTable, candidate.ToColumns[0])
		baseScore := candidate.Confidence
		evidenceItems := evidence.Evidence[candidateID]
		// Aggregate evidence scores
		for _, item := range evidenceItems {
			baseScore += item.Confidence * 0.2 // Weight evidence
		}
		// Normalize to 0-1
		if baseScore > 1.0 {
			baseScore = 1.0
		}
		scoring.Scores[candidateID] = baseScore
	}
	return scoring
}
func (p *ScanPipeline) reconcileGraph(metadata *MetadataResults, candidates []Candidate, scoring *ScoringResults) *GraphResults {
	graph := &GraphResults{
		Nodes: []GraphNode{},
		Edges: []GraphEdge{},
	}
	// Add nodes for all tables
	tableMap := make(map[string]bool)
	for _, table := range metadata.Tables {
		tableID := fmt.Sprintf("%s.%s", table.Schema, table.Name)
		if !tableMap[tableID] {
			graph.Nodes = append(graph.Nodes, GraphNode{
				ID:     tableID,
				Type:   "TABLE",
				Schema: table.Schema,
				Table:  table.Name,
			})
			tableMap[tableID] = true
		}
	}
	// Add edges for high-confidence candidates
	for _, candidate := range candidates {
		candidateID := fmt.Sprintf("%s.%s.%s->%s.%s.%s",
			candidate.FromSchema, candidate.FromTable, candidate.FromColumns[0],
			candidate.ToSchema, candidate.ToTable, candidate.ToColumns[0])
		score := scoring.Scores[candidateID]
		if score >= 0.5 { // Threshold for inclusion
			fromID := fmt.Sprintf("%s.%s", candidate.FromSchema, candidate.FromTable)
			toID := fmt.Sprintf("%s.%s", candidate.ToSchema, candidate.ToTable)
			graph.Edges = append(graph.Edges, GraphEdge{
				From:       fromID,
				To:         toID,
				Type:       "PK_FK",
				Confidence: score,
			})
		}
	}
	return graph
}
func (p *ScanPipeline) generateReport(metadata *MetadataResults, graph *GraphResults, scoring *ScoringResults) *ReportResults {
	report := &ReportResults{
		Summary: ReportSummary{
			TotalTables: len(metadata.Tables),
			TotalColumns: len(metadata.Columns),
			PKFKRelationships: len(graph.Edges),
		},
		Recommendations: []string{},
	}
	// Count confidence levels
	for _, score := range scoring.Scores {
		if score >= 0.8 {
			report.Summary.HighConfidence++
		} else if score >= 0.5 {
			report.Summary.MediumConfidence++
		} else {
			report.Summary.LowConfidence++
		}
	}
	// Generate recommendations
	if report.Summary.LowConfidence > 0 {
		report.Recommendations = append(report.Recommendations,
			fmt.Sprintf("Review %d low-confidence relationships", report.Summary.LowConfidence))
	}
	if report.Summary.PKFKRelationships == 0 {
		report.Recommendations = append(report.Recommendations,
			"No PK/FK relationships detected. Consider running with deeper profiling.")
	}
	return report
}
