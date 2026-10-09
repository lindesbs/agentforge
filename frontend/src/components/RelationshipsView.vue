<script setup lang="ts">
import type { Capability, ConfigFile, Definition, Project } from '../types'
defineProps<{ project: Project | null; definitions: Record<string, Definition>; capabilities: Capability[]; blocked: boolean }>()
defineEmits<{ open: [file: ConfigFile]; workspace: [] }>()
const providers = [{ id: 'codex', name: 'Codex' }, { id: 'claude', name: 'Claude Code' }]
</script>

<template>
 <section aria-labelledby="relationships-heading">
  <div class="page-heading"><div><span class="eyebrow">Read-only configuration structure</span><h1 id="relationships-heading" tabindex="-1">Definition map</h1><p>See which native definitions belong to each provider and which references are explicitly declared.</p></div><span class="badge">Structure, not execution</span></div>
  <div v-if="!project" class="empty-state panel"><h2>Open a project to see its definitions</h2><p>The map uses discovered paths and details of files you explicitly open. It does not scan other content.</p><button class="secondary" @click="$emit('workspace')">Go to project workspace</button></div>
  <template v-else>
   <div class="map-root"><span class="eyebrow">Project</span><h2>{{ project.name }}</h2><p class="project-path">{{ project.root }}</p></div>
   <p class="field-help map-legend">Connecting lines mean “contains”. References mean “declared in file”, not execution order, inheritance or permission to open another file.</p>
   <div class="definition-map">
    <section v-for="provider in providers" :key="provider.id" class="provider-branch" :aria-label="provider.name + ' definitions'">
     <div class="section-heading"><h2>{{ provider.name }}</h2><span class="count">{{ project.files.filter(file=>file.provider===provider.id).length }} {{ project.files.filter(file=>file.provider===provider.id).length === 1 ? 'file' : 'files' }}</span></div>
     <p v-if="!project.files.some(file=>file.provider===provider.id)" class="empty-copy">No supported files discovered for this provider.</p>
     <ul class="map-files">
      <li v-for="file in project.files.filter(file=>file.provider===provider.id)" :key="file.path" class="map-file">
       <div class="chips"><span class="badge">{{ file.kind }}</span><span class="muted">{{ definitions[file.path] ? 'Loaded snapshot' : 'Contents not opened' }}</span></div>
       <button class="map-file-link" :disabled="blocked" :aria-label="'Open definition ' + file.path" @click="$emit('open',file)">{{ file.path }}</button>
       <p v-if="definitions[file.path]?.issues.some(issue=>issue.severity==='error')" class="warning">Definition needs attention. Open it to inspect diagnostics.</p>
       <ul v-if="definitions[file.path]?.agents.length" class="declared-roles">
        <li v-for="(agent,index) in definitions[file.path].agents" :key="index"><span class="eyebrow">{{ file.provider === 'codex' ? 'Declared role' : 'Agent definition' }}</span><strong>{{ agent.name }}</strong><p v-if="agent.description">{{ agent.description }}</p><div v-if="agent.configFile" class="reference-node"><span class="muted">References · not opened or verified</span><code>{{ agent.configFile }}</code></div></li>
       </ul>
      </li>
     </ul>
    </section>
   </div>
   <section class="panel capability-panel" aria-labelledby="capability-heading">
    <h2 id="capability-heading">Provider support in AgentForge</h2><p class="muted">Compatibility is format-specific. No automatic conversion between providers, no agent execution.</p>
    <ul class="capability-list"><li v-for="item in capabilities" :key="item.provider + item.kind"><div><strong>{{ item.provider === 'codex' ? 'Codex' : 'Claude Code' }} · {{ item.kind }}</strong><span class="badge">{{ item.format }}</span></div><p>{{ item.validation }}</p><p class="field-help">Native editor{{ item.templateSupport ? ' · Same-provider templates' : '' }}{{ item.structuredFields.length ? ' · Simple fields: ' + item.structuredFields.join(', ') : ' · No structured field editor' }}</p></li></ul>
    <p class="field-help">Only documented syntax and selected fields are checked, not the complete vendor schema. Unknown settings stay unchanged.</p>
   </section>
  </template>
 </section>
</template>
