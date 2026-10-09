<script setup lang="ts">
import { ref } from 'vue'

const dialog = ref<HTMLDialogElement | null>(null)
const destination = ref('')
const keepButton = ref<HTMLButtonElement | null>(null)
const discardButton = ref<HTMLButtonElement | null>(null)
let settle: ((discard: boolean) => void) | null = null

function confirm(nextAction: string): Promise<boolean> {
  if (!dialog.value || settle) return Promise.resolve(false)
  destination.value = nextAction
  return new Promise(resolve => {
    settle = resolve
    dialog.value!.showModal()
  })
}

function finish(discard: boolean) {
  dialog.value?.close()
  settle?.(discard)
  settle = null
}

function trapFocus(event: KeyboardEvent) {
  if (event.shiftKey && document.activeElement === keepButton.value) {
    event.preventDefault()
    discardButton.value?.focus()
  } else if (!event.shiftKey && document.activeElement === discardButton.value) {
    event.preventDefault()
    keepButton.value?.focus()
  }
}

defineExpose({ confirm })
</script>

<template>
  <dialog ref="dialog" class="discard-dialog" aria-labelledby="discard-title" aria-describedby="discard-description" @cancel.prevent="finish(false)" @keydown.tab="trapFocus">
    <span class="eyebrow">Unsaved work</span>
    <h2 id="discard-title">Discard unsaved changes?</h2>
    <p id="discard-description">{{ destination }} will replace your current draft. Your saved file will stay unchanged.</p>
    <div class="actions dialog-actions">
      <button ref="keepButton" class="secondary" autofocus @click="finish(false)">Keep editing</button>
      <button ref="discardButton" class="danger" @click="finish(true)">Discard changes</button>
    </div>
  </dialog>
</template>
