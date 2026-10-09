<script setup lang="ts">
import { ref } from 'vue'

type ConfigFile = { path: string; provider: string; kind: string }
type Project = { root: string; name: string; frameworks: string[]; files: ConfigFile[] }
type Document = { path: string; content: string; hash: string }
type Preview = { path: string; diff: string; changed: boolean }
type Issue = { severity: string; message: string }
type AgentFields = { path: string; hash: string; name: string; description: string }
type Bridge = {
 InspectProject: (path: string) => Promise<Project>
 SelectProjectDirectory: () => Promise<string>
 ReadConfig: (root: string, path: string) => Promise<Document>
 PreviewConfig: (root: string, path: string, hash: string, content: string) => Promise<Preview>
 SaveConfig: (root: string, path: string, hash: string, content: string) => Promise<Document>
 ReadAgentFields: (root: string, path: string) => Promise<AgentFields>
 PrepareAgentFields: (root: string, path: string, hash: string, name: string, description: string) => Promise<Document>
 ValidateConfig: (path: string, content: string) => Promise<Issue[]>
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
const busy = ref(false)
const error = ref('')
const message = ref('')
function bridge(): Bridge { const api=window.go?.main?.App; if(!api) throw new Error('Open AgentForge as a desktop application.'); return api }
function clearEditor() { document.value=null; draft.value=''; preview.value=null; agentFields.value=null; issues.value=[]; message.value='' }
async function inspectSelected(api: Bridge) {
 clearEditor()
 project.value = null
 project.value = await api.InspectProject(path.value)
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
  preview.value=null;message.value='Saved. A recovery backup was created beside the original file.'
 } catch(e){error.value=String(e);preview.value=null}finally{busy.value=false}
}
function draftChanged(){preview.value=null;issues.value=[];message.value=''}
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
