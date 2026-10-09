import { test as base, expect, type Page } from '@playwright/test'

// UI contract fixture only. Production code never imports or installs a mock
// bridge. Filesystem and provider behavior are tested by the Go suites.
export async function installBridge(page: Page) {
  await page.addInitScript(() => {
    const state = window as any
    state.__calls = []
    state.__reads = []
    state.__failInspect = false
    state.__failSave = false
    state.__delaySave = 0
    state.__delayLearning = 0
    state.__failImport = false
    state.__delayImport = 0
    state.__templates = [{ id: 'template-1', name: 'Team defaults', provider: 'codex', kind: 'instructions', sourcePath: 'AGENTS.md', createdAt: '', tags: ['team','review'] }]
    state.__documents = {
      'AGENTS.md': '# Project instructions\nKeep changes focused.\n',
      '.codex/config.toml': 'model = "example-model"\nsandbox_mode = "workspace-write"\n',
      '.claude/agents/reviewer.md': '---\nname: reviewer\ndescription: Reviews project changes\nmodel: sonnet\n---\nReview carefully.\n',
    }
    const files = [
      { path: '.claude/agents/reviewer.md', provider: 'claude', kind: 'agent' },
      { path: '.codex/config.toml', provider: 'codex', kind: 'settings' },
      { path: 'AGENTS.md', provider: 'codex', kind: 'instructions' },
    ]
    const doc = (path: string) => ({ path, content: state.__documents[path], hash: 'fixture-hash' })
    const api = {
      async InspectProject(root: string) {
        if (state.__failInspect) throw new Error('Project directory is not accessible. Check the path and try again.')
        return { root, name: 'sample-project', frameworks: ['Go', 'Node.js'], files: root.endsWith('/empty') ? [] : files, diagnostics: [] }
      },
      async SelectProjectDirectory() { return '/projects/sample-project' },
      async ReadConfig(_root: string, path: string) { state.__reads.push(path); return doc(path) },
      async DescribeConfig(path: string) {
        return {
          provider: path.startsWith('.claude') ? 'claude' : 'codex', kind: path.endsWith('.toml') ? 'settings' : 'instructions',
          fields: path.endsWith('.toml') ? [{ name: 'model', value: 'example-model' }] : [], agents: path.endsWith('.toml') ? (state.__agents ?? []) : [],
          issues: path.endsWith('.toml') ? [{ code: 'limited-validation', severity: 'info', message: 'Only syntax and selected field types are checked.' }] : [],
        }
      },
      async ReadAgentFields(_root: string, path: string) {
        const content = state.__documents[path]
        return { path, hash: 'fixture-hash', name: content.match(/^name: (.+)$/m)?.[1] ?? '', description: content.match(/^description: (.+)$/m)?.[1] ?? '' }
      },
      async PrepareAgentFields(_root: string, path: string, hash: string, name: string, description: string) {
        return { path, hash, content: `---\nname: ${name}\ndescription: ${description}\n---\nReview carefully.\n` }
      },
      async ValidateConfig(_path: string, content: string) {
        return content.includes('INVALID') ? [{ severity: 'error', message: 'Correct the invalid configuration before saving.' }] : []
      },
      async PreviewConfig(_root: string, path: string, _hash: string, content: string) {
        state.__calls.push('preview')
        return { path, changed: content !== state.__documents[path], diff: `--- a/${path}\n+++ b/${path}\n-${state.__documents[path]}\n+${content}` }
      },
      async SaveConfig(_root: string, path: string, _hash: string, content: string) {
        state.__calls.push('save')
        if (state.__delaySave) await new Promise(resolve => setTimeout(resolve, state.__delaySave))
        if (state.__failSave) throw new Error('File changed on disk; reload before editing.')
        state.__documents[path] = content
        return doc(path)
      },
      async ProviderCapabilities() { return [
        {provider:'codex',kind:'settings',format:'TOML',structuredFields:[],validation:'Syntax and selected field types',templateSupport:true},
        {provider:'codex',kind:'instructions',format:'Markdown',structuredFields:[],validation:'Non-empty content advisory',templateSupport:true},
        {provider:'claude',kind:'agent',format:'Markdown + YAML',structuredFields:['name','description'],validation:'Frontmatter syntax and selected field types',templateSupport:true},
      ] },
      async CheckPortability() { return [{severity:'warning',message:'Review native content for credentials, local paths and project-specific settings.'}] },
      async ListTemplates() { return structuredClone(state.__templates) },
      async SaveTaggedTemplate(_root: string, path: string, name: string, tags: string[]) {
        const file = files.find(file=>file.path===path)!
        const item = {id:'template-'+(state.__templates.length+1),name,provider:file.provider,kind:file.kind,sourcePath:path,tags:tags.map(tag=>tag.trim().toLowerCase()).filter(Boolean),createdAt:''}
        state.__templates.push(item);state.__calls.push('save-template');return item
      },
      async ReviewTemplateImport(payload: string) {
        const envelope = JSON.parse(payload)
        if(envelope.version!==1 || envelope.format!=='agentforge-template') throw new Error('Unsupported template format or version')
        return {template:envelope.template,issues:envelope.template.content.includes('INVALID') ? [{severity:'error',message:'Invalid native configuration'}] : [{severity:'warning',message:'Check for credentials and local paths before importing.'}]}
      },
      async ImportTemplate(payload: string) {
        state.__calls.push('import-template')
        if(state.__delayImport) await new Promise(resolve=>setTimeout(resolve,state.__delayImport))
        if(state.__failImport) throw new Error('Library is not writable. Retry after checking permissions.')
        const item={...JSON.parse(payload).template,id:'import-'+state.__templates.length,createdAt:''};state.__templates.push(item);return item
      },
      async ExportTemplate(id: string) { const {id: _id, createdAt: _createdAt, ...item}=state.__templates.find((item: any)=>item.id===id);return JSON.stringify({format:'agentforge-template',version:1,template:{...item,content:'# Team defaults\nReview before saving.\n'}}) },
      async SaveProjectAsTemplate() { return { id: 'template-1' } },
      async PrepareTemplateApply(_id: string, _root: string, path: string) { return { path, hash: 'fixture-hash', content: '# Team defaults\nReview before saving.\n' } },
      async ListLearnings() {
        if (state.__delayLearning) await new Promise(resolve => setTimeout(resolve, state.__delayLearning))
        return []
      },
      async SaveLearning(entry: unknown) { state.__calls.push('save-learning'); return entry },
      async FormatLearning() { return '# Learning\nRule and check.' },
    }
    state.go = { main: { App: api } }
  })
}

export const test = base.extend({
  page: async ({ page }, use) => {
    await installBridge(page)
    await use(page)
  },
})
export { expect }
