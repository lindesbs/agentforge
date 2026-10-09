<script setup lang="ts">
import { ref } from 'vue'

type ConfigFile = { path: string; provider: string; kind: string }
type Project = { root: string; name: string; frameworks: string[]; files: ConfigFile[] }
type InspectorBridge = { InspectProject: (path: string) => Promise<Project> }
declare global { interface Window { go?: { main?: { App?: InspectorBridge } } } }
const path = ref('')
const project = ref<Project | null>(null)
const busy = ref(false)
const error = ref('')
async function inspect() {
 error.value = ''
 project.value = null
 if (!window.go?.main?.App) { error.value = 'Open this interface in the AgentForge desktop application.'; return }
 busy.value = true
 try { project.value = await window.go.main.App.InspectProject(path.value) }
 catch (e) { error.value = String(e) }
 finally { busy.value = false }
}
</script>

<template>
 <main>
  <header><div class="brand">◆ AgentForge</div><span class="tag">Configuration only · Read-only preview</span></header>
  <section class="hero">
   <h1>Project inspector</h1>
   <p>Inspect local Codex and Claude Code configuration files without running any agents.</p>
   <form @submit.prevent="inspect"><label for="root">Project directory</label><div class="input-row"><input id="root" v-model="path" placeholder="/home/user/projects/my-app" required /><button :disabled="busy">{{ busy ? 'Inspecting…' : 'Inspect project' }}</button></div></form>
   <p v-if="error" class="error" role="alert">{{ error }}</p>
  </section>
  <section v-if="project" class="results">
   <h2>{{ project.name }}</h2><p class="muted">{{ project.root }}</p>
   <div class="chips"><span v-for="item in project.frameworks" :key="item" class="chip">{{ item }}</span><span v-if="!project.frameworks.length" class="muted">No framework marker detected</span></div>
   <h3>Configuration files ({{ project.files.length }})</h3>
   <p v-if="!project.files.length" class="muted">No supported project-local configuration files found.</p>
   <ul><li v-for="file in project.files" :key="file.path"><code>{{ file.path }}</code><span class="muted">{{ file.provider }} · {{ file.kind }}</span></li></ul>
  </section>
 </main>
</template>
