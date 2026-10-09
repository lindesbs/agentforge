<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import type { Bridge, ConfigFile, ImportReview, TemplateSummary } from '../types'

const props = defineProps<{ active: boolean; api: () => Bridge; target: ConfigFile | null; blocked: boolean }>()
const emit = defineEmits<{ updated: []; apply: [id: string] }>()
const entries = ref<TemplateSummary[]>([])
const query = ref('')
const provider = ref('')
const loading = ref(false)
const working = ref(false)
const error = ref('')
const message = ref('')
const payload = ref('')
const reviewedPayload = ref('')
const review = ref<ImportReview | null>(null)
const exportText = ref('')
const exportName = ref('')
const reviewHeading = ref<HTMLElement | null>(null)
const exportField = ref<HTMLTextAreaElement | null>(null)
const importField = ref<HTMLTextAreaElement | null>(null)
const savedHeading = ref<HTMLElement | null>(null)
let exportTrigger: HTMLElement | null = null
const visible = computed(() => entries.value.filter(item => (!provider.value || item.provider === provider.value) && `${item.name} ${item.sourcePath} ${(item.tags ?? []).join(' ')}`.toLowerCase().includes(query.value.trim().toLowerCase())))
const disabled = computed(() => loading.value || working.value || props.blocked)
const canImport = computed(() => review.value && reviewedPayload.value === payload.value && !review.value.issues.some(issue => issue.severity === 'error'))

function compatible(item: TemplateSummary) { return props.target?.provider === item.provider && props.target?.kind === item.kind }
async function refresh() {
 if (disabled.value) return
 loading.value = true; error.value = ''
 try { entries.value = await props.api().ListTemplates() }
 catch (e) { error.value = String(e) }
 finally { loading.value = false }
}
watch(() => props.active, active => { if (active) void refresh() }, { immediate: true })
function invalidateReview() { review.value = null; reviewedPayload.value = ''; message.value = '' }
async function cancelReview() { invalidateReview(); await nextTick(); importField.value?.focus() }
async function closeExport() { exportText.value='';exportName.value='';await nextTick();exportTrigger?.focus() }
async function reviewImport() {
 if (disabled.value) return
 error.value = ''; message.value = ''; working.value = true; review.value = null
 try {
  const result = await props.api().ReviewTemplateImport(payload.value)
  review.value = result; reviewedPayload.value = payload.value
  await nextTick(); reviewHeading.value?.focus()
 } catch (e) { error.value = String(e) }
 finally { working.value = false }
}
async function importReviewed() {
 if (disabled.value || !canImport.value) return
 working.value = true; error.value = ''; message.value = ''
 let imported = false
 try {
  const item = await props.api().ImportTemplate(reviewedPayload.value)
  payload.value = ''; review.value = null; reviewedPayload.value = ''
  message.value = `Imported “${item.name}” as a new local template. No project files changed.`
  emit('updated')
  imported = true
 } catch (e) { error.value = String(e) }
 finally { working.value = false }
 if (imported) {await refresh();await nextTick();savedHeading.value?.focus()}
}
async function exportItem(item: TemplateSummary) {
 if (disabled.value) return
 exportTrigger = document.activeElement as HTMLElement | null
 working.value = true; error.value = ''; exportText.value = ''
 try {
  exportText.value = await props.api().ExportTemplate(item.id); exportName.value = item.name
  await nextTick(); exportField.value?.focus()
 } catch (e) { error.value = String(e) }
 finally { working.value = false }
}
</script>

<template>
 <section aria-labelledby="templates-heading">
  <div class="page-heading"><div><span class="eyebrow">Reusable native configuration</span><h1 id="templates-heading" tabindex="-1">Template library</h1><p>Find local templates, review their compatibility and copy them into an existing file.</p></div><span class="badge">Local · No synchronization</span></div>
  <p v-if="error" class="notice error" role="alert">{{ error }}</p>
  <p v-if="message" class="notice success" role="status">{{ message }}</p>
  <div class="panel library-browser" :aria-busy="loading || working">
   <div class="section-heading"><h2 ref="savedHeading" tabindex="-1">Saved templates</h2><span class="count">{{ visible.length }}</span></div>
   <div class="library-filters"><div><label for="library-search">Search templates</label><input id="library-search" v-model="query" type="search" placeholder="Name, tag or source path" /></div><div><label for="library-provider">Provider</label><select id="library-provider" v-model="provider"><option value="">All providers</option><option value="codex">Codex</option><option value="claude">Claude Code</option></select></div><button class="secondary" :disabled="disabled" @click="refresh">Refresh templates</button></div>
   <p class="field-help">Create a tagged template from an opened file using Configuration templates in the project editor.</p>
   <p class="notice" :class="target ? '' : 'warning'">{{ target ? 'Current target: ' + target.path + '. Only matching provider and file types can be loaded.' : 'No target selected. Open a project file to enable loading compatible templates.' }}</p>
   <p v-if="loading" role="status">Loading templates…</p>
   <div v-else-if="!error && !visible.length" class="empty-state compact"><h3>{{ entries.length ? 'No matching templates' : 'Your reusable configurations belong here' }}</h3><p>{{ entries.length ? 'Try another name, tag or provider.' : 'Save a project file as a template, or review an exported template below.' }}</p><button v-if="entries.length" class="secondary small" @click="query='';provider=''">Clear template filters</button></div>
   <ul class="template-cards">
    <li v-for="item in visible" :key="item.id" class="template-card">
     <div class="section-heading"><h3>{{ item.name }}</h3><span class="badge" :class="compatible(item) ? 'badge-neutral' : ''">{{ compatible(item) ? 'Compatible with target' : target ? 'Different target type' : 'Choose a target' }}</span></div>
     <p class="muted">{{ item.provider === 'codex' ? 'Codex' : 'Claude Code' }} · {{ item.kind }}</p><p class="project-path">{{ item.sourcePath }}</p>
     <div v-if="item.tags?.length" class="chips template-tags"><span v-for="tag in item.tags" :key="tag" class="badge">{{ tag }}</span></div>
     <div class="actions"><button class="secondary small" :disabled="disabled || !compatible(item)" @click="emit('apply',item.id)">Load into current file</button><button class="secondary small" :disabled="disabled" :aria-label="'Export ' + item.name" @click="exportItem(item)">Review export</button></div>
    </li>
   </ul>
   <p class="field-help warning">Loading replaces the complete draft; it does not merge project overrides. Review the diff, then save explicitly.</p>
  </div>

  <section v-if="exportText" class="panel transfer-panel" aria-labelledby="export-heading">
   <h2 id="export-heading">Export: {{ exportName }}</h2><p class="notice warning">This export includes the complete native text. Check for credentials and private paths before copying or sharing it.</p>
   <label for="template-export">Portable template JSON</label><textarea id="template-export" ref="exportField" class="code-text" :value="exportText" readonly rows="10" spellcheck="false"></textarea><p class="field-help">Select and copy the JSON into a local .json file. Referenced files are not bundled. Nothing is uploaded.</p>
   <button class="secondary small" @click="closeExport">Close export</button>
  </section>

  <section class="panel transfer-panel" aria-labelledby="import-heading">
   <h2 id="import-heading">Import a template</h2><p class="muted">Paste AgentForge template JSON. Review is read-only; importing creates a new library entry, never a project file.</p>
   <form @submit.prevent="reviewImport"><fieldset :disabled="disabled"><legend class="sr-only">Template import</legend><label for="template-import">Template JSON to review</label><textarea id="template-import" ref="importField" v-model="payload" rows="6" class="code-text" spellcheck="false" required maxlength="6299648" @input="invalidateReview"></textarea><button class="secondary" type="submit">{{ working ? 'Working…' : 'Review import' }}</button></fieldset></form>
   <div v-if="review" class="import-review">
    <h3 ref="reviewHeading" tabindex="-1">Review: {{ review.template.name }}</h3><p>{{ review.template.provider }} · {{ review.template.kind }} · {{ review.template.sourcePath }}</p>
    <div class="diagnostics"><ul><li v-for="(issue,index) in review.issues" :key="index" :class="issue.severity"><strong>{{ issue.severity.toUpperCase() }}</strong> · {{ issue.message }}</li></ul></div>
    <label for="import-native">Native content to import</label><textarea id="import-native" :value="review.template.content" readonly class="code-text" rows="7"></textarea>
    <p class="field-help">Tags: {{ review.template.tags.join(', ') || 'None' }}. Existing templates are never overwritten.</p>
    <div class="actions"><button :disabled="disabled || !canImport" @click="importReviewed">Import reviewed template</button><button class="secondary" :disabled="disabled" @click="cancelReview">Cancel review</button></div>
   </div>
  </section>
 </section>
</template>
