<script setup lang="ts">
import { onMounted, ref } from 'vue'

type ConfigFile = { path: string; provider: string; kind: string }
type Project = { root: string; name: string; frameworks: string[]; files: ConfigFile[] }
type Document = { path: string; content: string; hash: string }
type Preview = { path: string; diff: string; changed: boolean }
type Issue = { severity: string; message: string }
type TemplateSummary = { id: string; name: string; provider: string; kind: string; sourcePath: string; createdAt: string }
type LearningEntry = { id: string; category: string; language: string; title: string; finding: string; rule: string; check: string; createdAt: string }
type LearningProjectPreview = { path: string; beforeHash: string; proposed: string; diff: string; changed: boolean }
type LearningImportCandidate = {category:string; index:number; sourcePath:string; sourceHash:string; entry:LearningEntry; duplicate:boolean}
type LearningImportSkipped = {sourcePath:string; heading:string; reason:string}
type LearningImportScan = {candidates:LearningImportCandidate[]; skipped:LearningImportSkipped[]}
type AgentFields = { path: string; hash: string; name: string; description: string }
type Bridge = {
 InspectProject: (path: string) => Promise<Project>
 StartupProject: () => Promise<string>
 SelectProjectDirectory: () => Promise<string>
 ReadConfig: (root: string, path: string) => Promise<Document>
 PreviewConfig: (root: string, path: string, hash: string, content: string) => Promise<Preview>
 SaveConfig: (root: string, path: string, hash: string, content: string) => Promise<Document>
 ReadAgentFields: (root: string, path: string) => Promise<AgentFields>
 PrepareAgentFields: (root: string, path: string, hash: string, name: string, description: string) => Promise<Document>
 ValidateConfig: (path: string, content: string) => Promise<Issue[]>
 ListTemplates: () => Promise<TemplateSummary[]>
 SaveProjectAsTemplate: (root: string, path: string, name: string) => Promise<TemplateSummary>
 PrepareTemplateApply: (id: string, root: string, targetPath: string) => Promise<Document>
 ListLearnings: () => Promise<LearningEntry[]>
 SaveLearning: (entry: LearningEntry) => Promise<LearningEntry>
 FormatLearning: (entry: LearningEntry) => Promise<string>
 PreviewLearningInProject: (root: string, id: string) => Promise<LearningProjectPreview>
 ApplyLearningToProject: (root: string, id: string, expectedHash: string, expectedProposed: string) => Promise<LearningProjectPreview>
 ScanProjectLearnings: (root: string) => Promise<LearningImportScan>
 ImportProjectLearning: (root: string, category: string, index: number, sourceHash: string, language: string) => Promise<LearningEntry>
}
declare global { interface Window { go?: { main?: { App?: Bridge } } } }
const path = ref('')
const project = ref<Project | null>(null)
const document = ref<Document | null>(null)
const draft = ref('')
const preview = ref<Preview | null>(null)
const agentFields = ref<AgentFields | null>(null)
const issues = ref<Issue[]>([])
const agentName = ref('')
const agentDescription = ref('')
const templates = ref<TemplateSummary[]>([])
const templateName = ref('')
const selectedTemplate = ref('')
const pendingTemplate = ref('')
const learningEntries = ref<LearningEntry[]>([])
const learningCategory = ref('ALLGEMEIN_OK')
const learningLanguage = ref('Allgemein')
const languageFilter = ref('Alle')
const learningTitle = ref('')
const learningFinding = ref('')
const learningRule = ref('')
const learningCheck = ref('')
const learningMarkdown = ref('')
const projectLearningPreview = ref<LearningProjectPreview | null>(null)
const projectLearningID = ref('')
const projectLearningRoot = ref('')
const learningImportScan = ref<LearningImportScan | null>(null)
const learningImportBusy = ref(false)
const learningImportLanguages = ref<Record<string,string>>({})
const learningBusy = ref(false)
const languageOptions = ['Allgemein','Go','PHP','TypeScript','JavaScript','Python','Rust','Kotlin','Java','C#','C++','Dart','SQL','Shell','Andere']
const categoryOptions = ['PROJEKT_OK','PROJEKT_NOK','ALLGEMEIN_OK','ALLGEMEIN_NOK']
const busy = ref(false)
const error = ref('')
const message = ref('')
onMounted(async () => {
 void refreshLearnings()
 try {
  const selected = await bridge().StartupProject()
  if(selected) {
   path.value=selected
   await inspect()
  }
 } catch(e) { error.value=String(e) }
})
async function refreshLearnings() {
 try { learningEntries.value = await bridge().ListLearnings() }
 catch(e) { error.value=String(e) }
}
async function saveLearning() {
 error.value='';message.value='';learningBusy.value=true
 try {
  const entry:LearningEntry = {id:'',createdAt:'',category:learningCategory.value,language:learningLanguage.value,
   title:learningTitle.value,finding:learningFinding.value,rule:learningRule.value,check:learningCheck.value}
  await bridge().SaveLearning(entry)
  await refreshLearnings()
  learningTitle.value='';learningFinding.value='';learningRule.value='';learningCheck.value=''
  message.value='Bestätigte Erkenntnis in der lokalen Vorlagenbibliothek gespeichert; Projektdateien unverändert.'
 } catch(e){error.value=String(e)}finally{learningBusy.value=false}
}
function importKey(candidate:LearningImportCandidate):string {
 return candidate.category+':'+candidate.index
}
async function scanExistingLearnings() {
 if(!project.value)return
 error.value='';message.value='';learningImportBusy.value=true
 try {
  learningImportScan.value=await bridge().ScanProjectLearnings(project.value.root)
  learningImportLanguages.value={}
  for(const c of learningImportScan.value.candidates) {
   learningImportLanguages.value[importKey(c)]=c.entry.language
  }
 }catch(e){error.value=String(e);learningImportScan.value=null}finally{learningImportBusy.value=false}
}
async function importExistingLearning(candidate:LearningImportCandidate) {
 if(!project.value)return
 error.value='';message.value='';learningImportBusy.value=true
 try {
  const chosen=learningImportLanguages.value[importKey(candidate)]||candidate.entry.language
  await bridge().ImportProjectLearning(project.value.root,candidate.category,candidate.index,candidate.sourceHash,chosen)
  await refreshLearnings()
  await scanExistingLearnings()
  message.value='Erkenntnis in die lokale Bibliothek importiert. Projektdateien blieben unverändert.'
 }catch(e){error.value=String(e)}finally{learningImportBusy.value=false}
}
async function previewLearningInProject(entry:LearningEntry) {
 if(!project.value)return
 error.value='';message.value='';projectLearningPreview.value=null;learningBusy.value=true
 try {
  const result=await bridge().PreviewLearningInProject(project.value.root,entry.id)
  projectLearningPreview.value=result;projectLearningID.value=entry.id;projectLearningRoot.value=project.value.root
 } catch(e){error.value=String(e)}finally{learningBusy.value=false}
}
async function applyLearningToProject() {
 if(!project.value||!projectLearningPreview.value||!projectLearningID.value)return
 error.value='';message.value='';learningBusy.value=true
 try {
  if(project.value.root!==projectLearningRoot.value)throw new Error('Projekt wurde gewechselt. Vorschau erneut erstellen.')
  const plan=projectLearningPreview.value
  await bridge().ApplyLearningToProject(project.value.root,projectLearningID.value,plan.beforeHash,plan.proposed)
  projectLearningPreview.value=null;projectLearningID.value=''
  message.value='Erkenntnis im Projektdokument ergänzt. Vorhandene Dokumente werden zuvor gesichert.'
 }catch(e){error.value=String(e);projectLearningPreview.value=null}finally{learningBusy.value=false}
}
async function showLearning(entry:LearningEntry) {
 error.value=''
 try {learningMarkdown.value=await bridge().FormatLearning(entry)}
 catch(e){error.value=String(e)}
}
function filteredLearnings():LearningEntry[] {
 return learningEntries.value.filter(e=>languageFilter.value==='Alle'||e.language===languageFilter.value)
}
async function refreshTemplates() {
 try { templates.value=await bridge().ListTemplates() }
 catch(e) {error.value=String(e)}
}
async function saveTemplate() {
 if(!project.value||!document.value||!templateName.value.trim())return
 error.value='';message.value='';busy.value=true
 try {
  // Snapshot the current on-disk file, not unsaved textarea changes.
  const saved=await bridge().SaveProjectAsTemplate(project.value.root,document.value.path,templateName.value.trim())
  await refreshTemplates()
  selectedTemplate.value=saved.id
  templateName.value=''
  message.value='Template saved from the current file on disk. Unsaved changes were not included.'
 } catch(e){error.value=String(e)}finally{busy.value=false}
}
async function applyTemplate() {
 if(!project.value||!document.value||!selectedTemplate.value)return
 error.value='';message.value='';preview.value=null;issues.value=[];busy.value=true
 try {
  const proposed=await bridge().PrepareTemplateApply(selectedTemplate.value,project.value.root,document.value.path)
  if(proposed.hash!==document.value.hash)throw new Error('Target changed. Reload the file before applying a template.')
  draft.value=proposed.content
  issues.value=await bridge().ValidateConfig(document.value.path,draft.value)
  preview.value=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  pendingTemplate.value=selectedTemplate.value
  message.value='Template applied to the editor only. Review the diff before saving.'
 } catch(e){error.value=String(e)}finally{busy.value=false}
}
function availableTemplates(): TemplateSummary[] {
 if(!document.value)return []
 const p=document.value.path
 const provider=p==='AGENTS.md'||p==='.codex/config.toml'?'codex':'claude'
 const kind=p==='AGENTS.md'||p==='CLAUDE.md'?'instructions':p.endsWith('.json')||p.endsWith('.toml')?'settings':'agent'
 return templates.value.filter(item=>item.provider===provider&&item.kind===kind)
}
function bridge(): Bridge { const api=window.go?.main?.App; if(!api) throw new Error('Open AgentForge as a desktop application.'); return api }
function clearEditor() { document.value=null; draft.value=''; preview.value=null; agentFields.value=null; issues.value=[]; selectedTemplate.value=''; pendingTemplate.value=''; message.value='' }
async function inspectSelected(api: Bridge) {
 clearEditor()
 projectLearningPreview.value=null
 projectLearningID.value=''
 learningImportScan.value=null
 project.value = null
 project.value = await api.InspectProject(path.value)
 await refreshTemplates()
 await refreshLearnings()
}
async function inspect() {
 error.value='';busy.value=true
 try { await inspectSelected(bridge()) } catch(e) { error.value=String(e) } finally { busy.value=false }
}
async function chooseDirectory() {
 error.value='';busy.value=true
 try {
  const api=bridge();const selected=await api.SelectProjectDirectory()
  if(selected) {path.value=selected;await inspectSelected(api)}
 } catch(e) {error.value=String(e)} finally {busy.value=false}
}
async function openFile(file: ConfigFile) {
 if(!project.value)return
 error.value='';clearEditor();busy.value=true
 try {
  const api=bridge();const doc=await api.ReadConfig(project.value.root,file.path);document.value=doc;draft.value=doc.content
  await refreshTemplates()
  if(file.provider==='claude'&&file.kind==='agent') {
   try {const fields=await api.ReadAgentFields(project.value.root,file.path);agentFields.value=fields;agentName.value=fields.name;agentDescription.value=fields.description}
   catch { agentFields.value=null }
  }
 }
 catch(e){error.value=String(e)}finally{busy.value=false}
}
async function applyStructuredFields() {
 if(!project.value||!document.value||!agentFields.value)return
 error.value='';message.value='';preview.value=null;busy.value=true
 try {
  const prepared=await bridge().PrepareAgentFields(project.value.root,document.value.path,document.value.hash,agentName.value,agentDescription.value)
  draft.value=prepared.content
  issues.value=await bridge().ValidateConfig(document.value.path,draft.value)
  preview.value=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
 } catch(e){error.value=String(e)} finally {busy.value=false}
}
async function previewChanges() {
 if(!project.value||!document.value)return
 error.value='';message.value='';preview.value=null;busy.value=true
 try {
  issues.value=await bridge().ValidateConfig(document.value.path,draft.value)
  preview.value=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
 }
 catch(e){error.value=String(e)}finally{busy.value=false}
}
async function saveChanges() {
 if(!project.value||!document.value||!preview.value?.changed)return
 error.value='';busy.value=true
 try {
  // Re-preview immediately before save, and require a matching draft.
  const current=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  if(!current.changed||current.diff!==preview.value.diff)throw new Error('Preview changed. Review again before saving.')
  document.value=await bridge().SaveConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  preview.value=null;pendingTemplate.value='';message.value='Saved. A recovery backup was created beside the original file.'
 } catch(e){error.value=String(e);preview.value=null}finally{busy.value=false}
}
function draftChanged(){preview.value=null;issues.value=[];pendingTemplate.value='';message.value=''}
</script>

<template>
 <main>
  <header><div class="brand">◆ AgentForge</div><span class="tag">Local configuration editor · No agent execution</span></header>
  <section class="hero">
   <h1>Project configuration</h1>
   <p>Inspect local Codex and Claude Code files. Edit their native text with a mandatory diff preview.</p>
   <form @submit.prevent="inspect">
    <label for="root">Project directory</label>
    <div class="input-row">
     <input id="root" v-model="path" placeholder="/home/user/projects/my-app" required />
     <button type="button" class="secondary" :disabled="busy" @click="chooseDirectory">Browse…</button>
     <button type="submit" :disabled="busy">{{ busy ? 'Working…' : 'Inspect project' }}</button>
    </div>
   </form>
   <p v-if="error" class="error" role="alert">{{ error }}</p>
   <p v-if="message" class="success" role="status">{{ message }}</p>
  </section>
  <section class="learning-panel">
   <h2>Entwicklungs-Erkenntnisse</h2>
   <p>Bestätigte Lösungen und Fehlerursachen als wiederverwendbare Vorlagen je Programmiersprache dokumentieren. Die vier Projektdateien werden hierbei nicht automatisch verändert.</p>
   <div class="learning-grid">
    <div><label for="learning-category">Geltungsbereich / Bewertung</label>
     <select id="learning-category" v-model="learningCategory"><option v-for="c in categoryOptions" :key="c">{{ c }}</option></select></div>
    <div><label for="learning-language">Programmiersprache</label>
     <select id="learning-language" v-model="learningLanguage"><option v-for="l in languageOptions" :key="l">{{ l }}</option></select></div>
   </div>
   <label for="learning-title">Titel der bestätigten Erkenntnis</label><input id="learning-title" v-model="learningTitle" maxlength="120" placeholder="Kurz und eindeutig" />
   <label for="learning-finding">Bestätigte Erkenntnis / Ursache</label><textarea id="learning-finding" v-model="learningFinding" rows="3" />
   <label for="learning-rule">Verbindliche Regel</label><textarea id="learning-rule" v-model="learningRule" rows="3" />
   <label for="learning-check">Überprüfbare Kontrolle</label><textarea id="learning-check" v-model="learningCheck" rows="2" placeholder="Regressionstest, Lint, Healthcheck …" />
   <div class="actions"><button :disabled="learningBusy || !learningTitle.trim() || !learningFinding.trim() || !learningRule.trim() || !learningCheck.trim()" @click="saveLearning">Erkenntnis als Vorlage speichern</button>
    <button class="secondary" :disabled="learningBusy" @click="refreshLearnings">Vorlagen aktualisieren</button></div>
   <section v-if="project" class="learning-import">
    <h3>Vorhandene Erkenntnisse aus dem Projekt importieren</h3>
    <p>Die vier `docs/PROJEKT_*.md`- und `docs/ALLGEMEIN_*.md`-Dateien werden nur gelesen. Einträge mit bestätigter Erkenntnis, Regel und Prüfung werden als Vorschau angeboten.</p>
    <button class="secondary" :disabled="learningImportBusy" @click="scanExistingLearnings">Projektdokumente durchsuchen</button>
    <template v-if="learningImportScan">
     <p>{{ learningImportScan.candidates.length }} importierbare Einträge · {{ learningImportScan.skipped.length }} unvollständige Abschnitte</p>
     <ul class="learning-import-list">
      <li v-for="candidate in learningImportScan.candidates" :key="importKey(candidate)">
       <div><strong>{{ candidate.entry.title }}</strong> <span class="muted">{{ candidate.sourcePath }}</span></div>
       <p><strong>Erkenntnis:</strong> {{ candidate.entry.finding }}</p>
       <p><strong>Regel:</strong> {{ candidate.entry.rule }}</p>
       <p><strong>Prüfung:</strong> {{ candidate.entry.check }}</p>
       <div class="actions">
        <select v-model="learningImportLanguages[importKey(candidate)]" :aria-label="'Sprache für '+candidate.entry.title">
         <option v-for="lang in languageOptions" :key="lang">{{ lang }}</option>
        </select>
        <button :disabled="learningImportBusy || candidate.duplicate && learningImportLanguages[importKey(candidate)] === candidate.entry.language" @click="importExistingLearning(candidate)">
         {{ candidate.duplicate && learningImportLanguages[importKey(candidate)] === candidate.entry.language ? 'Bereits vorhanden' : 'In Bibliothek übernehmen' }}
        </button>
       </div>
      </li>
     </ul>
     <details v-if="learningImportScan.skipped.length">
      <summary>Nicht importierte Abschnitte anzeigen</summary>
      <ul><li v-for="(item,i) in learningImportScan.skipped" :key="i">{{ item.sourcePath }} – {{ item.heading }}: {{ item.reason }}</li></ul>
     </details>
    </template>
   </section>
   <div class="learning-grid"><h3>Gespeicherte Erkenntnisse ({{ filteredLearnings().length }})</h3>
    <select v-model="languageFilter" aria-label="Sprache filtern"><option>Alle</option><option v-for="l in languageOptions" :key="l">{{ l }}</option></select></div>
   <ul class="learning-list"><li v-for="item in filteredLearnings()" :key="item.id">
    <button class="file-button" @click="showLearning(item)">{{ item.title }}</button>
    <span class="muted">{{ item.language }} · {{ item.category }}</span>
    <button v-if="project" class="secondary learning-apply-button" :disabled="learningBusy" @click="previewLearningInProject(item)">In Projekt vorschlagen</button>
   </li></ul>
   <section v-if="projectLearningPreview" class="learning-apply-preview">
    <h3>Vorschau: {{ projectLearningPreview.path }}</h3>
    <p>Diese Änderung wird erst nach ausdrücklicher Bestätigung in das ausgewählte Projekt geschrieben.</p>
    <pre>{{ projectLearningPreview.diff }}</pre>
    <div class="actions">
     <button :disabled="learningBusy || !projectLearningPreview.changed" @click="applyLearningToProject">Geprüfte Erkenntnis übernehmen</button>
     <button class="secondary" :disabled="learningBusy" @click="projectLearningPreview = null">Abbrechen</button>
    </div>
   </section>
   <div v-if="learningMarkdown" class="learning-output"><label for="learning-markdown">Markdown-Vorlage zum Übernehmen in die passende docs/*_OK/NOK.md</label>
    <textarea id="learning-markdown" readonly :value="learningMarkdown" rows="9"></textarea>
   </div>
  </section>
  <section v-if="project" class="results">
   <h2>{{ project.name }}</h2><p class="muted">{{ project.root }}</p>
   <div class="chips"><span v-for="item in project.frameworks" :key="item" class="chip">{{ item }}</span><span v-if="!project.frameworks.length" class="muted">No framework marker detected</span></div>
   <h3>Configuration files ({{ project.files.length }})</h3>
   <p v-if="!project.files.length" class="muted">No supported project-local configuration files found.</p>
   <ul><li v-for="file in project.files" :key="file.path"><button class="file-button" :disabled="busy" @click="openFile(file)">{{ file.path }}</button><span class="muted">{{ file.provider }} · {{ file.kind }}</span></li></ul>
   <div v-if="document" class="editor">
    <h3>Edit: {{ document.path }}</h3>
    <div v-if="agentFields" class="structured-fields">
     <h4>Claude agent properties</h4>
     <label for="agent-name">Name</label><input id="agent-name" v-model="agentName" />
     <label for="agent-description">Description</label><input id="agent-description" v-model="agentDescription" />
     <button :disabled="busy" @click="applyStructuredFields">Apply fields and preview</button>
     <p class="muted">Only existing simple YAML name/description fields are supported. Other content is preserved.</p>
    </div>
    <section class="template-library">
      <h4>Global template library</h4>
      <p class="muted">Store a copy of the current file or load a compatible template into the editor. Nothing is automatically synchronized or written to your project.</p>
      <div class="input-row">
       <input v-model="templateName" aria-label="New template name" placeholder="Template name" maxlength="120" />
       <button class="secondary" :disabled="busy || !templateName.trim()" @click="saveTemplate">Save current file as template</button>
      </div>
      <div class="input-row">
       <select v-model="selectedTemplate" aria-label="Select template">
        <option value="">Choose compatible template…</option>
        <option v-for="item in availableTemplates()" :key="item.id" :value="item.id">{{ item.name }} · {{ item.sourcePath }}</option>
       </select>
       <button class="secondary" :disabled="busy || !selectedTemplate" @click="applyTemplate">Preview template on this file</button>
      </div>
      <p class="muted">{{ templates.length }} saved template(s) · {{ availableTemplates().length }} compatible with this file</p>
    </section>
    <p class="muted">Raw native configuration text · 1 MiB limit · Existing files only</p>
    <textarea v-model="draft" spellcheck="false" rows="16" @input="draftChanged"></textarea>
    <div class="actions"><button :disabled="busy || draft===document.content" @click="previewChanges">Preview changes</button><button class="secondary" :disabled="busy" @click="openFile({path:document.path,provider:'',kind:''})">Reload file</button></div>
    <div v-if="issues.length" class="diagnostics" aria-live="polite"><h4>Configuration diagnostics</h4><ul><li v-for="(issue,index) in issues" :key="index" :class="issue.severity">{{ issue.severity.toUpperCase() }}: {{ issue.message }}</li></ul></div>
    <div v-if="preview" class="preview">
     <h3>Change preview</h3><p v-if="!preview.changed">No changes detected.</p>
     <template v-else><pre>{{ preview.diff }}</pre><button :disabled="busy || issues.some(i => i.severity === 'error')" @click="saveChanges">Save reviewed changes</button></template>
    </div>
   </div>
  </section>
 </main>
</template>
