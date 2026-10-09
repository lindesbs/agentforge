export type ConfigFile = { path: string; provider: string; kind: string }
export type DiscoveryDiagnostic = { path: string; code: string; severity: string; message: string }
export type Project = { root: string; name: string; frameworks: string[]; files: ConfigFile[]; diagnostics: DiscoveryDiagnostic[] }
export type Document = { path: string; content: string; hash: string }
export type Preview = { path: string; diff: string; changed: boolean }
export type Issue = { severity: string; message: string }
export type TemplateSummary = { id: string; name: string; provider: string; kind: string; sourcePath: string; createdAt: string; tags: string[] }
export type LearningEntry = { id: string; category: string; language: string; title: string; finding: string; rule: string; check: string; createdAt: string }
export type AgentFields = { path: string; hash: string; name: string; description: string }
export type Definition = {
 provider: string; kind: string
 fields: { name: string; value: string }[]
 agents: { name: string; description: string; configFile: string }[]
 issues: { code: string; severity: string; message: string }[]
}
export type Bridge = {
 InspectProject: (path: string) => Promise<Project>
 SelectProjectDirectory: () => Promise<string>
 ReadConfig: (root: string, path: string) => Promise<Document>
 PreviewConfig: (root: string, path: string, hash: string, content: string) => Promise<Preview>
 SaveConfig: (root: string, path: string, hash: string, content: string) => Promise<Document>
 ReadAgentFields: (root: string, path: string) => Promise<AgentFields>
 PrepareAgentFields: (root: string, path: string, hash: string, name: string, description: string) => Promise<Document>
 ValidateConfig: (path: string, content: string) => Promise<Issue[]>
 DescribeConfig: (path: string, content: string) => Promise<Definition>
 ProviderCapabilities: () => Promise<Capability[]>
 CheckPortability: (path: string, content: string) => Promise<Issue[]>
 SaveTaggedTemplate: (root: string, path: string, name: string, tags: string[]) => Promise<TemplateSummary>
 ReviewTemplateImport: (payload: string) => Promise<ImportReview>
 ImportTemplate: (payload: string) => Promise<TemplateSummary>
 ExportTemplate: (id: string) => Promise<string>
 ListTemplates: () => Promise<TemplateSummary[]>
 SaveProjectAsTemplate: (root: string, path: string, name: string) => Promise<TemplateSummary>
 PrepareTemplateApply: (id: string, root: string, targetPath: string) => Promise<Document>
 ListLearnings: () => Promise<LearningEntry[]>
 SaveLearning: (entry: LearningEntry) => Promise<LearningEntry>
 FormatLearning: (entry: LearningEntry) => Promise<string>
}
export type Capability = { provider: string; kind: string; format: string; structuredFields: string[]; validation: string; templateSupport: boolean }
export type PortableTemplate = { name: string; provider: string; kind: string; sourcePath: string; content: string; tags: string[] }
export type ImportReview = { template: PortableTemplate; issues: Issue[] }
