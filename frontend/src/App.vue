<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import DiscardDialog from './components/DiscardDialog.vue'

import type { ConfigFile, Project, Document, Preview, Issue, TemplateSummary, LearningEntry, AgentFields, Definition, Bridge, Capability } from './types'
import TemplateLibrary from './components/TemplateLibrary.vue'
import RelationshipsView from './components/RelationshipsView.vue'

declare global { interface Window { go?: { main?: { App?: Bridge } } } }
const path = ref('')
const project = ref<Project | null>(null)
const document = ref<Document | null>(null)
const definition = ref<Definition | null>(null)
const draft = ref('')
const preview = ref<Preview | null>(null)
const agentFields = ref<AgentFields | null>(null)
const issues = ref<Issue[]>([])
const agentName = ref('')
const agentDescription = ref('')
const templates = ref<TemplateSummary[]>([])
const templateName = ref('')
const templateTags = ref('')
const templateQuery = ref('')
const capabilities = ref<Capability[]>([])
const loadedDefinitions = ref<Record<string, Definition>>({})
const portability = ref<Issue[]>([])
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
const learningBusy = ref(false)
const learningLoading = ref(false)
const languageOptions = ['Allgemein','Go','PHP','TypeScript','JavaScript','Python','Rust','Kotlin','Java','C#','C++','Dart','SQL','Shell','Andere']
const categoryOptions = ['PROJEKT_OK','PROJEKT_NOK','ALLGEMEIN_OK','ALLGEMEIN_NOK']
const categoryLabels: Record<string, string> = {
 PROJEKT_OK: 'Project · Solution', PROJEKT_NOK: 'Project · Prevention',
 ALLGEMEIN_OK: 'General · Solution', ALLGEMEIN_NOK: 'General · Prevention',
}
function languageLabel(value: string) { return value === 'Allgemein' ? 'General' : value === 'Andere' ? 'Other' : value }
const busy = ref(false)
const operation = ref('')
const error = ref('')
const message = ref('')
const projectError = ref('')
const learningError = ref('')
const learningMessage = ref('')
type View = 'workspace' | 'learning' | 'templates' | 'relationships'
const view = ref<View>('workspace')
const fileQuery = ref('')
const editorHeading = ref<HTMLElement | null>(null)
const previewHeading = ref<HTMLElement | null>(null)
const learningHeading = ref<HTMLElement | null>(null)
const workspaceHeading = ref<HTMLElement | null>(null)
const discardDialog = ref<InstanceType<typeof DiscardDialog> | null>(null)
const rawDirty = computed(() => !!document.value && draft.value !== document.value.content)
const fieldsDirty = computed(() => !!agentFields.value && (agentName.value !== agentFields.value.name || agentDescription.value !== agentFields.value.description))
const dirty = computed(() => rawDirty.value || fieldsDirty.value)
const visibleFiles = computed(() => project.value?.files.filter(file => `${file.path} ${file.provider}`.toLowerCase().includes(fileQuery.value.trim().toLowerCase())) ?? [])
const diffLines = computed(() => preview.value?.diff.split('\n') ?? [])

function beforeUnload(event: BeforeUnloadEvent) {
 if (dirty.value) { event.preventDefault(); event.returnValue = '' }
}
onMounted(() => window.addEventListener('beforeunload', beforeUnload))
onUnmounted(() => window.removeEventListener('beforeunload', beforeUnload))
async function allowDiscard(action: string, needed = dirty.value) {
 return !needed || await discardDialog.value?.confirm(action) === true
}
async function changeView(next: View) {
 view.value = next
 await nextTick()
 const heading = next === 'workspace' ? workspaceHeading.value : next === 'learning' ? learningHeading.value : window.document.getElementById(next + '-heading')
 heading?.focus()
 if (next === 'learning' && !learningBusy.value) await refreshLearnings()
}
async function focusPreview() { await nextTick(); previewHeading.value?.focus() }
function fileName(path: string) { return path.split('/').pop() ?? path }
function fileFolder(path: string) { return path.includes('/') ? path.slice(0, path.lastIndexOf('/')) : 'Project root' }
function diffClass(line: string) {
 if (line.startsWith('---') || line.startsWith('+++')) return 'diff-file'
 return line.startsWith('+') ? 'diff-added' : line.startsWith('-') ? 'diff-removed' : ''
}
async function refreshLearnings() {
 if(learningLoading.value)return
 learningError.value='';learningLoading.value=true
 try { learningEntries.value = await bridge().ListLearnings() }
 catch(e) { learningError.value=String(e) }
 finally { learningLoading.value=false }
}
async function saveLearning() {
 if (learningBusy.value || learningLoading.value) return
 learningError.value='';learningMessage.value='';learningBusy.value=true
 try {
  const entry:LearningEntry = {id:'',createdAt:'',category:learningCategory.value,language:learningLanguage.value,
   title:learningTitle.value,finding:learningFinding.value,rule:learningRule.value,check:learningCheck.value}
  await bridge().SaveLearning(entry)
  await refreshLearnings()
  learningTitle.value='';learningFinding.value='';learningRule.value='';learningCheck.value=''
  learningMessage.value='Learning saved to your local library. No project files were changed.'
 } catch(e){learningError.value=String(e)}finally{learningBusy.value=false}
}
async function showLearning(entry:LearningEntry) {
 learningError.value=''
 try {learningMarkdown.value=await bridge().FormatLearning(entry)}
 catch(e){learningError.value=String(e)}
}
function filteredLearnings():LearningEntry[] {
 return learningEntries.value.filter(e=>languageFilter.value==='Alle'||e.language===languageFilter.value)
}
async function refreshTemplates() {
 try { templates.value=await bridge().ListTemplates() }
 catch(e) {error.value=String(e)}
}
async function saveTemplate() {
 if(busy.value||!project.value||!document.value||!templateName.value.trim())return
 error.value='';message.value='';busy.value=true
 try {
  // Snapshot the current on-disk file, not unsaved textarea changes.
  const saved=await bridge().SaveTaggedTemplate(project.value.root,document.value.path,templateName.value.trim(),templateTags.value.split(','))
  await refreshTemplates()
  selectedTemplate.value=saved.id
  templateName.value=''
  templateTags.value=''
  message.value='Template saved from the current file on disk. Unsaved changes were not included.'
 } catch(e){error.value=String(e)}finally{busy.value=false}
}
async function applyTemplate() {
 if(busy.value||!project.value||!document.value||!selectedTemplate.value)return
 if(!await allowDiscard('Loading this template'))return
 error.value='';message.value='';preview.value=null;issues.value=[];busy.value=true
 try {
  const proposed=await bridge().PrepareTemplateApply(selectedTemplate.value,project.value.root,document.value.path)
  if(proposed.hash!==document.value.hash)throw new Error('Target changed. Reload the file before applying a template.')
  draft.value=proposed.content
  await loadDetails(document.value)
  issues.value=await bridge().ValidateConfig(document.value.path,draft.value)
  preview.value=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  pendingTemplate.value=selectedTemplate.value
  message.value='Template applied to the editor only. It replaces the whole file, including project-specific overrides; nothing is merged automatically. Review every removed line before saving.'
  await focusPreview()
  return true
 } catch(e){error.value=String(e)}finally{busy.value=false}
}
function availableTemplates(): TemplateSummary[] {
 if(!document.value)return []
 const p=document.value.path
 const provider=p==='AGENTS.md'||p==='.codex/config.toml'?'codex':'claude'
 const kind=p==='AGENTS.md'||p==='CLAUDE.md'?'instructions':p.endsWith('.json')||p.endsWith('.toml')?'settings':'agent'
 return templates.value.filter(item=>item.provider===provider&&item.kind===kind&&`${item.name} ${item.sourcePath} ${(item.tags ?? []).join(' ')}`.toLowerCase().includes(templateQuery.value.trim().toLowerCase()))
}
async function useLibraryTemplate(id: string) {
 selectedTemplate.value=id
 if(await applyTemplate()) {await changeView('workspace');await focusPreview()}
}
async function openRelated(file: ConfigFile) {
 if(await openFile(file)) {await changeView('workspace');await nextTick();editorHeading.value?.focus()}
}
const currentCapability = computed(() => capabilities.value.find(item=>item.provider===definition.value?.provider&&item.kind===definition.value?.kind))
function bridge(): Bridge { const api=window.go?.main?.App; if(!api) throw new Error('Open AgentForge as a desktop application.'); return api }
function clearEditor() { document.value=null; definition.value=null; draft.value=''; preview.value=null; agentFields.value=null; issues.value=[]; selectedTemplate.value=''; pendingTemplate.value=''; templateName.value='';templateTags.value='';templateQuery.value='';portability.value=[]; message.value='';error.value='' }
async function inspectSelected(api: Bridge) {
 const nextProject = await api.InspectProject(path.value)
 clearEditor()
 project.value = nextProject
 loadedDefinitions.value={}
 fileQuery.value=''
 await refreshTemplates()
 capabilities.value=await api.ProviderCapabilities()
}
async function inspect() {
 if(busy.value||!await allowDiscard('Opening this project'))return
 projectError.value='';busy.value=true;operation.value='project'
 try { await inspectSelected(bridge()) } catch(e) { projectError.value=String(e) } finally { busy.value=false;operation.value='' }
}
async function chooseDirectory() {
 if(busy.value||!await allowDiscard('Choosing another project'))return
 projectError.value='';busy.value=true;operation.value='project'
 try {
  const api=bridge();const selected=await api.SelectProjectDirectory()
  if(selected) {path.value=selected;await inspectSelected(api)}
 } catch(e) {projectError.value=String(e)} finally {busy.value=false;operation.value=''}
}
async function loadDetails(doc: Document) {
 definition.value=null;agentFields.value=null
 definition.value=await bridge().DescribeConfig(doc.path,doc.content)
 loadedDefinitions.value={...loadedDefinitions.value,[doc.path]:definition.value}
 portability.value=await bridge().CheckPortability(doc.path,doc.content)
 if(project.value && /^\.claude\/agents\/[^/]+\.md$/i.test(doc.path)) {
  try {
   const fields=await bridge().ReadAgentFields(project.value.root,doc.path)
   agentFields.value=fields;agentName.value=fields.name;agentDescription.value=fields.description
  } catch { agentFields.value=null }
 }
}
async function openFile(file: ConfigFile) {
 if(busy.value||!project.value||!await allowDiscard('Opening this file'))return
 error.value='';busy.value=true;operation.value='file'
 try {
  const doc=await bridge().ReadConfig(project.value.root,file.path)
  clearEditor();document.value=doc;draft.value=doc.content
  await loadDetails(doc)
  await refreshTemplates()
  await nextTick();editorHeading.value?.focus()
  return true
 }
 catch(e){error.value=String(e)}finally{busy.value=false;operation.value=''}
}
async function applyStructuredFields() {
 if(busy.value||!project.value||!document.value||!agentFields.value)return
 if(!await allowDiscard('Applying these fields',rawDirty.value))return
 error.value='';message.value='';preview.value=null;busy.value=true
 try {
  const prepared=await bridge().PrepareAgentFields(project.value.root,document.value.path,document.value.hash,agentName.value,agentDescription.value)
  draft.value=prepared.content
  agentFields.value={...agentFields.value,name:agentName.value,description:agentDescription.value}
  issues.value=await bridge().ValidateConfig(document.value.path,draft.value)
  preview.value=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  await focusPreview()
 } catch(e){error.value=String(e)} finally {busy.value=false}
}
async function previewChanges() {
 if(busy.value||!project.value||!document.value||fieldsDirty.value)return
 error.value='';message.value='';preview.value=null;busy.value=true;operation.value='preview'
 try {
  issues.value=await bridge().ValidateConfig(document.value.path,draft.value)
  preview.value=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  await focusPreview()
 }
 catch(e){error.value=String(e)}finally{busy.value=false;operation.value=''}
}
async function saveChanges() {
 if(busy.value||!project.value||!document.value||fieldsDirty.value||!preview.value?.changed||issues.value.some(i=>i.severity==='error'))return
 error.value='';busy.value=true;operation.value='save'
 try {
  // Re-preview immediately before save, and require a matching draft.
  const current=await bridge().PreviewConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  if(!current.changed||current.diff!==preview.value.diff)throw new Error('Preview changed. Review again before saving.')
  document.value=await bridge().SaveConfig(project.value.root,document.value.path,document.value.hash,draft.value)
  preview.value=null;pendingTemplate.value='';message.value='Saved. A recovery backup was created beside the original file.'
  await loadDetails(document.value)
  await nextTick();editorHeading.value?.focus()
 } catch(e){error.value=String(e);preview.value=null}finally{busy.value=false;operation.value=''}
}
function draftChanged(){preview.value=null;issues.value=[];pendingTemplate.value='';message.value=''}
</script>

<template>
 <div class="app-shell">
  <a class="skip-link" href="#main-content">Skip to workspace</a>
  <header class="app-header">
   <div class="brand"><span class="brand-mark" aria-hidden="true">AF</span><div>AgentForge<span class="brand-caption">Native configuration workspace</span></div></div>
   <span class="local-status"><span aria-hidden="true" class="status-dot"></span>Local files · No agent execution</span>
  </header>
  <div class="app-body">
   <nav class="app-nav" aria-label="Main navigation">
    <p class="nav-label">Workspace</p>
    <button :aria-current="view === 'workspace' ? 'page' : undefined" :disabled="busy" @click="changeView('workspace')">
     <span aria-hidden="true" class="nav-symbol">▤</span>Project configuration<span v-if="dirty" class="nav-dirty" aria-label="Unsaved changes">●</span>
    </button>
    <button :aria-current="view === 'learning' ? 'page' : undefined" :disabled="busy" @click="changeView('learning')">
     <span aria-hidden="true" class="nav-symbol">▧</span>Learning library
    </button>
    <button :aria-current="view === 'templates' ? 'page' : undefined" :disabled="busy" @click="changeView('templates')"><span aria-hidden="true" class="nav-symbol">▱</span>Template library</button>
    <button :aria-current="view === 'relationships' ? 'page' : undefined" :disabled="busy" @click="changeView('relationships')"><span aria-hidden="true" class="nav-symbol">⌘</span>Definition map</button>
    <div class="nav-context">
     <span class="eyebrow">Current project</span>
     <strong>{{ project?.name ?? 'No project open' }}</strong>
     <p>{{ project ? 'Configuration stays in your project. Review every change before saving.' : 'Open a project to discover its Codex and Claude Code files.' }}</p>
    </div>
    <p class="nav-footer">Inspect. Edit. Review.</p>
   </nav>
   <main id="main-content" tabindex="-1">
    <p v-if="(view === 'templates' || view === 'relationships') && error" class="notice error" role="alert">{{ error }}</p>
    <div v-show="view === 'workspace'" class="workspace-view">
     <div class="page-heading">
      <div><span class="eyebrow">Project workspace</span><h1 ref="workspaceHeading" tabindex="-1">Project configuration</h1><p>Keep native agent files in your control, from inspection to reviewed changes.</p></div>
      <span class="badge">Codex + Claude Code</span>
     </div>
     <form class="project-picker" :aria-busy="operation === 'project'" @submit.prevent="inspect">
      <label for="root">Project directory</label>
      <div class="input-row">
       <input id="root" v-model="path" :disabled="busy" placeholder="/home/user/projects/my-app" required spellcheck="false" aria-describedby="project-help" />
       <button type="button" class="secondary" :disabled="busy" @click="chooseDirectory">Browse…</button>
       <button type="submit" :disabled="busy">{{ operation === 'project' ? 'Inspecting…' : 'Inspect project' }}</button>
      </div>
      <p id="project-help" class="field-help">Discovery reads file names. Open a file to inspect its contents.</p>
      <p v-if="projectError" class="notice error" role="alert">{{ projectError }}</p>
     </form>

     <section v-if="!project" class="empty-state welcome-state" aria-labelledby="welcome-title">
      <span class="empty-symbol" aria-hidden="true">⌘</span>
      <h2 id="welcome-title">Your project, its native configuration.</h2>
      <p>Choose a local directory to find instruction files, agent definitions and provider settings.</p>
      <ol class="workflow-steps"><li><span>01</span>Open a project</li><li><span>02</span>Choose a file</li><li><span>03</span>Review before saving</li></ol>
      <p class="muted">No sign-in. No model calls. No agents running in the background.</p>
     </section>

     <section v-else class="project-content" aria-label="Project files and editor">
      <div class="project-heading">
       <div><h2>{{ project.name }}</h2><p class="project-path">{{ project.root }}</p></div>
       <div class="chips"><span v-for="item in project.frameworks" :key="item" class="badge">{{ item }}</span><span v-if="!project.frameworks.length" class="muted">No framework marker detected</span></div>
      </div>
      <details v-if="project.diagnostics.length" class="disclosure discovery-diagnostics">
       <summary>Discovery notices <span class="count">{{ project.diagnostics.length }}</span></summary>
       <div class="diagnostics"><ul><li v-for="item in project.diagnostics" :key="item.path + ':' + item.code" :class="item.severity"><strong>{{ item.severity.toUpperCase() }}</strong> · <code>{{ item.path }}</code>: {{ item.message }}</li></ul></div>
      </details>
      <div class="editor-layout">
       <nav class="file-explorer" aria-label="Configuration files">
        <div class="section-heading"><h3>Configuration files</h3><span class="count">{{ project.files.length }}</span></div>
        <label class="field-label" for="file-search">Filter files</label>
        <input id="file-search" v-model="fileQuery" type="search" placeholder="File name or provider" />
        <p v-if="!project.files.length" class="empty-copy">No supported files found. Add an AGENTS.md or a provider configuration to this project, then inspect it again.</p>
        <div v-else-if="!visibleFiles.length" class="empty-copy"><p>No files match “{{ fileQuery }}”.</p><button class="secondary small" @click="fileQuery = ''">Clear filter</button></div>
        <ul v-else class="file-list">
         <li v-for="file in visibleFiles" :key="file.path">
          <button class="file-item" :aria-current="document?.path === file.path ? 'true' : undefined" :disabled="busy" :aria-label="file.path" @click="openFile(file)">
           <span class="file-item-title"><span class="file-glyph" aria-hidden="true">{{ file.provider === 'codex' ? 'C' : 'A' }}</span><strong>{{ fileName(file.path) }}</strong></span>
           <span class="file-item-path">{{ fileFolder(file.path) }}</span>
           <span class="file-item-meta">{{ file.provider === 'codex' ? 'Codex' : 'Claude Code' }} · {{ file.kind }}</span>
          </button>
         </li>
        </ul>
       </nav>

       <section class="editor-panel" :aria-busy="busy" aria-label="Configuration editor">
        <div v-if="!document" class="empty-state file-empty">
         <span class="empty-symbol" aria-hidden="true">≡</span><h2>Choose a configuration file</h2><p>Inspect its definition, edit the native text and review a diff before anything is saved.</p>
         <p v-if="operation === 'file'" role="status">Opening file…</p>
         <p v-if="error" class="notice error" role="alert">{{ error }}</p>
        </div>
        <template v-else>
         <div class="editor-heading">
          <div><span class="eyebrow">Native file editor</span><h2 ref="editorHeading" tabindex="-1">{{ fileName(document.path) }}</h2><p class="project-path">{{ document.path }}</p></div>
          <span class="badge" :class="dirty ? 'badge-warning' : 'badge-neutral'">{{ dirty ? 'Unsaved changes' : 'Matches loaded file' }}</span>
         </div>
         <p v-if="error" class="notice error" role="alert">{{ error }}</p>
         <p v-if="message" class="notice success" role="status">{{ message }}</p>

         <details v-if="definition" class="disclosure definition-summary">
          <summary>Definition &amp; provider checks <span v-if="definition.issues.some(item => item.severity === 'error')" class="badge badge-warning">Needs attention</span><span v-else class="summary-meta">{{ definition.provider }} · {{ definition.kind }}</span></summary>
          <div class="disclosure-content">
           <p class="muted">Read-only details from the loaded file. Unsaved draft changes are not reflected here.</p>
           <div v-if="currentCapability" class="capability-note"><strong>{{ currentCapability.format }} · Native text editing</strong><p>{{ currentCapability.validation }}. {{ currentCapability.structuredFields.length ? 'Structured fields: ' + currentCapability.structuredFields.join(', ') + ' (simple scalars only).' : 'No structured field editor for this format.' }}</p><p>Unknown settings and comments are retained. Complete vendor-schema compatibility is not guaranteed.</p></div>
           <dl v-if="definition.fields.length" class="definition-fields"><template v-for="field in definition.fields" :key="field.name"><dt>{{ field.name }}</dt><dd>{{ field.value || '(empty)' }}</dd></template></dl>
           <template v-if="definition.provider === 'codex' && definition.agents.length">
            <h3>Declared agent roles</h3>
            <ul class="role-list"><li v-for="agent in definition.agents" :key="agent.name"><strong>{{ agent.name }}</strong><p v-if="agent.description">{{ agent.description }}</p><p v-if="agent.configFile" class="muted">Referenced configuration: <code>{{ agent.configFile }}</code> (not opened)</p></li></ul>
           </template>
           <p v-if="!definition.fields.length && !definition.agents.length" class="muted">No supported structured fields to display.</p>
           <div v-if="definition.issues.length" class="diagnostics"><h3>Provider inspection diagnostics</h3><ul><li v-for="(issue,index) in definition.issues" :key="index" :class="issue.severity"><strong>{{ issue.severity.toUpperCase() }}</strong> · {{ issue.message }}</li></ul></div>
           <div v-if="portability.some(item=>item.severity==='warning')" class="diagnostics"><h3>Portability checks</h3><ul><li v-for="(issue,index) in portability.filter(item=>item.severity==='warning')" :key="index" class="warning">{{ issue.message }}</li></ul></div>
          </div>
         </details>

         <fieldset class="editor-fields" :disabled="busy">
          <legend class="sr-only">Edit native configuration</legend>
          <details v-if="agentFields" class="disclosure structured-fields">
           <summary>Simple Claude agent fields <span v-if="fieldsDirty" class="badge badge-warning">Not applied</span></summary>
           <div class="disclosure-content">
            <label for="agent-name">Name</label><input id="agent-name" v-model="agentName" @input="draftChanged" />
            <label for="agent-description">Description</label><input id="agent-description" v-model="agentDescription" @input="draftChanged" />
            <p class="field-help">Existing single-line fields only. Applying fields creates a draft and a diff.</p>
            <button class="secondary" :disabled="!fieldsDirty" @click="applyStructuredFields">Apply fields and preview</button>
           </div>
          </details>
          <div class="source-heading"><label for="native-content">Native configuration</label><span class="muted">UTF-8 · Up to 1 MiB</span></div>
          <textarea id="native-content" v-model="draft" class="source-editor" spellcheck="false" rows="18" aria-describedby="editor-help" @input="draftChanged"></textarea>
          <p id="editor-help" class="field-help">Changes stay in this draft until you review the diff and explicitly save.</p>
          <p v-if="fieldsDirty" class="field-help warning">Apply your structured fields to the draft before previewing.</p>
         </fieldset>
         <div class="editor-actions">
          <button :disabled="busy || !rawDirty || fieldsDirty" @click="previewChanges">{{ operation === 'preview' ? 'Preparing diff…' : 'Preview changes' }}</button>
          <button class="secondary" :disabled="busy" @click="openFile({path:document.path,provider:'',kind:''})">Reload file</button>
          <span v-if="busy" class="muted" role="status">{{ operation === 'save' ? 'Saving reviewed changes…' : 'Working…' }}</span>
         </div>

         <div v-if="issues.length" class="diagnostics draft-diagnostics" role="status"><h3>Draft validation</h3><ul><li v-for="(issue,index) in issues" :key="index" :class="issue.severity"><strong>{{ issue.severity.toUpperCase() }}</strong> · {{ issue.message }}</li></ul></div>
         <section v-if="preview" class="preview" aria-label="Change preview">
          <div class="section-heading"><h3 ref="previewHeading" tabindex="-1">Review your changes</h3><span class="badge">{{ preview.changed ? 'Not saved yet' : 'No changes' }}</span></div>
          <p v-if="!preview.changed">The draft matches the current file.</p>
          <template v-else>
           <p class="field-help">Minus lines are removed; plus lines are added. Saving creates a recovery backup beside the original.</p>
           <pre class="diff" aria-label="Changes to the native file"><code><span v-for="(line,index) in diffLines" :key="index" class="diff-line" :class="diffClass(line)">{{ line || ' ' }}
</span></code></pre>
           <p v-if="issues.some(i => i.severity === 'error')" class="notice error">Resolve the draft validation errors and preview again before saving.</p>
           <div class="actions"><button :disabled="busy || issues.some(i => i.severity === 'error')" @click="saveChanges">{{ operation === 'save' ? 'Saving…' : 'Save reviewed changes' }}</button><span class="muted">Only this file will be changed.</span></div>
          </template>
         </section>

         <details class="disclosure template-library">
          <summary>Configuration templates <span class="count">{{ templates.length }}</span></summary>
          <div class="disclosure-content">
           <p class="muted">Copy a saved file into your local library, or load a compatible template into this draft.</p>
           <fieldset :disabled="busy"><legend class="sr-only">Template actions</legend>
            <label for="template-name">New template name</label>
            <div class="input-row"><input id="template-name" v-model="templateName" placeholder="For example, team review guidelines" maxlength="120" /><button class="secondary" :disabled="!templateName.trim()" @click="saveTemplate">Save file as template</button></div>
            <label for="template-tags">Tags (optional, comma-separated)</label><input id="template-tags" v-model="templateTags" placeholder="team, review, go" maxlength="395" />
            <p class="field-help">Uses the file on disk, not your unsaved draft. Review its contents for sensitive settings first.</p>
            <label for="template-select">Compatible template</label>
            <label for="template-search" class="sr-only">Search compatible templates</label><input id="template-search" v-model="templateQuery" type="search" placeholder="Search names or tags" />
            <div class="input-row"><select id="template-select" v-model="selectedTemplate"><option value="">Choose a template…</option><option v-for="item in availableTemplates()" :key="item.id" :value="item.id">{{ item.name }} · {{ item.sourcePath }}</option></select><button class="secondary" :disabled="!selectedTemplate" @click="applyTemplate">Load and preview</button></div>
            <p class="field-help">{{ availableTemplates().length }} compatible template(s). Loading never saves the project file.</p>
            <p class="field-help warning">Whole-file replacement, not a merge. Project-specific overrides can be removed; inspect the diff before saving.</p>
            <p v-if="pendingTemplate" class="field-help">A template is loaded in this draft. Review the diff before saving.</p>
           </fieldset>
          </div>
         </details>
        </template>
       </section>
      </div>
     </section>
    </div>

    <TemplateLibrary v-show="view === 'templates'" :active="view === 'templates'" :api="bridge" :target="project?.files.find(file=>file.path===document?.path) ?? null" :blocked="busy" @updated="refreshTemplates" @apply="useLibraryTemplate" />
    <RelationshipsView v-show="view === 'relationships'" :project="project" :definitions="loadedDefinitions" :capabilities="capabilities" :blocked="busy" @open="openRelated" @workspace="changeView('workspace')" />
    <section v-show="view === 'learning'" class="learning-view" aria-labelledby="learning-title-heading">
     <div class="page-heading"><div><span class="eyebrow">Local knowledge library</span><h1 id="learning-title-heading" ref="learningHeading" tabindex="-1">Learning library</h1><p>Turn confirmed findings into reusable rules, with a check that proves they work.</p></div><span class="badge">Project files stay unchanged</span></div>
     <p v-if="learningError" class="notice error" role="alert">{{ learningError }}</p>
     <p v-if="learningMessage" class="notice success" role="status">{{ learningMessage }}</p>
     <div class="learning-layout">
      <form class="learning-form panel" @submit.prevent="saveLearning">
       <h2>Record a confirmed learning</h2><p class="muted">All fields are required. Store causes and checks, never credentials.</p>
       <fieldset :disabled="learningBusy"><legend class="sr-only">Learning details</legend>
        <div class="form-columns">
         <div><label for="learning-category">Scope &amp; outcome</label><select id="learning-category" v-model="learningCategory"><option v-for="c in categoryOptions" :key="c" :value="c">{{ categoryLabels[c] }}</option></select></div>
         <div><label for="learning-language">Language</label><select id="learning-language" v-model="learningLanguage"><option v-for="l in languageOptions" :key="l" :value="l">{{ languageLabel(l) }}</option></select></div>
        </div>
        <label for="learning-title">Learning title</label><input id="learning-title" v-model="learningTitle" maxlength="120" required placeholder="A specific, reusable finding" />
        <label for="learning-finding">Confirmed finding or cause</label><textarea id="learning-finding" v-model="learningFinding" rows="3" required></textarea>
        <label for="learning-rule">Rule to follow</label><textarea id="learning-rule" v-model="learningRule" rows="3" required></textarea>
        <label for="learning-check">How to verify it</label><textarea id="learning-check" v-model="learningCheck" rows="2" required placeholder="A regression test, check or observable result"></textarea>
        <div class="actions"><button type="submit" :disabled="learningLoading">{{ learningBusy ? 'Saving learning…' : 'Save learning' }}</button></div>
       </fieldset>
      </form>
      <div class="learning-saved panel">
       <div class="section-heading"><h2>Saved learnings</h2><span class="count">{{ filteredLearnings().length }}</span></div>
       <label for="language-filter">Filter by language</label><select id="language-filter" v-model="languageFilter"><option value="Alle">All languages</option><option v-for="l in languageOptions" :key="l" :value="l">{{ languageLabel(l) }}</option></select>
       <p v-if="learningLoading" class="notice" role="status">Loading saved learnings…</p>
       <div v-else-if="!learningError && !filteredLearnings().length" class="empty-state compact"><h3>{{ learningEntries.length ? 'No matching learnings' : 'Build your reference library' }}</h3><p>{{ learningEntries.length ? 'Choose another language to see your saved entries.' : 'Your confirmed findings will appear here. Select one to get its Markdown template.' }}</p></div>
       <ul class="learning-list"><li v-for="item in filteredLearnings()" :key="item.id"><button class="learning-item" @click="showLearning(item)"><strong>{{ item.title }}</strong><span>{{ languageLabel(item.language) }} · {{ categoryLabels[item.category] ?? item.category }}</span></button></li></ul>
       <button class="secondary small" :disabled="learningBusy || learningLoading" @click="refreshLearnings">Refresh library</button>
       <div v-if="learningMarkdown" class="learning-output"><label for="learning-markdown">Markdown to copy into your project</label><textarea id="learning-markdown" class="code-text" readonly :value="learningMarkdown" rows="9" aria-describedby="markdown-help"></textarea><p id="markdown-help" class="field-help">Copy into the appropriate docs/*_OK.md or docs/*_NOK.md file after review.</p></div>
      </div>
     </div>
    </section>
   </main>
  </div>
  <DiscardDialog ref="discardDialog" />
 </div>
</template>
