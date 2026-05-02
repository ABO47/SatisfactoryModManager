<script lang="ts">
  import T from '$lib/components/T.svelte';
  import { customGamePaths, addCustomGamePath, removeCustomGamePath } from '$lib/store/settingsStore';
  import { OpenDirectoryDialog } from '$wailsjs/go/app/app';

  export let parent: { onClose: () => void };

  let customPathError: string | null = null;
  let fileDialogOpen = false;
  let isAdding = false;

  async function pickCustomGamePath() {
    if(fileDialogOpen) {
      return;
    }
    fileDialogOpen = true;
    try {
      let result = await OpenDirectoryDialog({
        defaultDirectory: '',
      });
      if (result) {
        await addCustomPathToList(result);
      }
    } catch (e) {
      if(e instanceof Error) {
        customPathError = e.message;
      } else if (typeof e === 'string') {
        customPathError = e;
      } else {
        customPathError = 'Unknown error';
      }
    } finally {
      fileDialogOpen = false;
    }
  }

  async function addCustomPathToList(path: string) {
    try {
      isAdding = true;
      customPathError = null;
      await addCustomGamePath(path);
    } catch(e) {
      if (e instanceof Error) {
        customPathError = e.message;
      } else if (typeof e === 'string') {
        customPathError = e;
      } else {
        customPathError = 'Unknown error';
      }
    } finally {
      isAdding = false;
    }
  }

  async function removeCustomPath(path: string) {
    try {
      customPathError = null;
      await removeCustomGamePath(path);
    } catch(e) {
      if (e instanceof Error) {
        customPathError = e.message;
      } else if (typeof e === 'string') {
        customPathError = e;
      } else {
        customPathError = 'Unknown error';
      }
    }
  }
</script>

<div style="max-height: calc(100vh - 3rem); max-width: calc(100vw - 3rem);" class="w-[60rem] card flex flex-col gap-2">
  <header class="card-header font-bold text-2xl text-center">
    <T defaultValue="Custom Game Locations" keyName="settings.customPaths.title" />
  </header>
  <section class="p-4 grow overflow-y-auto">
    <div class="flex flex-col gap-4">
      <div>
        <button 
          class="btn btn-primary"
          on:click={() => pickCustomGamePath()}
          disabled={isAdding || fileDialogOpen}
        >
          <T defaultValue="Add Custom Location" keyName="settings.customPaths.add" />
        </button>
        {#if customPathError}
          <p class="text-error mt-2">{customPathError}</p>
        {/if}
      </div>

      {#if $customGamePaths.length > 0}
        <div class="divider my-0"></div>
        <div class="space-y-2">
          {#each $customGamePaths as path (path)}
            <div class="flex items-center justify-between p-2 bg-base-200 rounded">
              <code class="text-sm break-all">{path}</code>
              <button 
                class="btn btn-sm btn-error btn-outline"
                on:click={() => removeCustomPath(path)}
              >
                <T defaultValue="Remove" keyName="settings.customPaths.remove" />
              </button>
            </div>
          {/each}
        </div>
      {:else}
        <p class="text-center text-base-content/50">
          <T defaultValue="No custom locations added" keyName="settings.customPaths.empty" />
        </p>
      {/if}
    </div>
  </section>
  <footer class="card-footer">
    <button class="btn btn-primary" on:click={() => parent.onClose()}>
      <T defaultValue="Close" keyName="common.close" />
    </button>
  </footer>
</div>
