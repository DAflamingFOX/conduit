<script lang="ts">
  // Svelte 5 Runes
  let activeTab = $state<'nodes' | 'storage' | 'trigger' | 'logs'>('nodes');
  let nodes = $state<any[]>([]);
  let storagePools = $state<Record<string, any>>({});
  let targetFilePath = $state('/tmp/sample.mp4');
  let triggerMessage = $state('');
  let logLines = $state<string[]>([
    '[SYSTEM] Conduit Manager Svelte 5 UI connected.',
    '[SYSTEM] Monitoring storage pools and registered node manifests.',
  ]);

  async function fetchNodes() {
    try {
      const res = await fetch('/api/v1/nodes');
      if (res.ok) {
        nodes = await res.json();
        addLog(`[API] Loaded ${nodes.length} node manifests from Manager.`);
      }
    } catch (e) {
      addLog(`[ERROR] Failed to fetch nodes: ${e}`);
    }
  }

  async function fetchStorage() {
    try {
      const res = await fetch('/api/v1/storage');
      if (res.ok) {
        storagePools = await res.json();
        addLog(`[API] Refreshed storage pool health probes.`);
      }
    } catch (e) {
      addLog(`[ERROR] Failed to fetch storage pools: ${e}`);
    }
  }

  async function handleTriggerJob(e: Event) {
    e.preventDefault();
    triggerMessage = 'Submitting job...';
    try {
      const res = await fetch('/api/v1/jobs/trigger', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ file_path: targetFilePath }),
      });
      const data = await res.json();
      if (res.ok) {
        triggerMessage = `[SUCCESS] Job Triggered! Job ID: ${data.job_id}`;
        addLog(`[JOB] Triggered execution for file: ${targetFilePath} (ID: ${data.job_id})`);
      } else {
        triggerMessage = `[ERROR] ${data.error}`;
      }
    } catch (err) {
      triggerMessage = `[ERROR] ${err}`;
    }
  }

  function addLog(msg: string) {
    const time = new Date().toISOString().split('T')[1].slice(0, 8);
    logLines = [...logLines, `[${time}] ${msg}`];
  }

  // Fetch data on load
  $effect(() => {
    fetchNodes();
    fetchStorage();
  });
</script>

<!-- Plain, Utilitarian Slate Layout -->
<div
  class="min-h-screen bg-slate-900 text-slate-200 font-mono text-sm p-4 border-t-2 border-slate-600"
>
  <!-- Header -->
  <header class="border-b border-slate-700 pb-3 mb-4 flex justify-between items-center">
    <div>
      <h1 class="text-base font-bold text-slate-100 uppercase tracking-wider">
        CONDUIT // FILE PROCESSING MANAGER
      </h1>
      <p class="text-xs text-slate-400">Utilitarian Open-Source Node Engine</p>
    </div>
    <div class="text-xs border border-slate-700 px-2 py-1 bg-slate-800 text-slate-300">
      STATUS: ONLINE | STACK: SVELTE 5 + GO
    </div>
  </header>

  <!-- Navigation Tabs -->
  <nav class="flex gap-2 mb-4 text-xs font-bold">
    <button
      onclick={() => (activeTab = 'nodes')}
      class="px-3 py-1.5 border border-slate-700 uppercase tracking-wide cursor-pointer text-left font-mono {activeTab ===
      'nodes'
        ? 'bg-slate-700 text-white'
        : 'bg-slate-800 text-slate-400 hover:text-slate-200'}"
    >
      [1] Node Registry ({nodes.length})
    </button>
    <button
      onclick={() => (activeTab = 'storage')}
      class="px-3 py-1.5 border border-slate-700 uppercase tracking-wide cursor-pointer text-left font-mono {activeTab ===
      'storage'
        ? 'bg-slate-700 text-white'
        : 'bg-slate-800 text-slate-400 hover:text-slate-200'}"
    >
      [2] Storage Pools ({Object.keys(storagePools).length})
    </button>
    <button
      onclick={() => (activeTab = 'trigger')}
      class="px-3 py-1.5 border border-slate-700 uppercase tracking-wide cursor-pointer text-left font-mono {activeTab ===
      'trigger'
        ? 'bg-slate-700 text-white'
        : 'bg-slate-800 text-slate-400 hover:text-slate-200'}"
    >
      [3] Trigger Job
    </button>
    <button
      onclick={() => (activeTab = 'logs')}
      class="px-3 py-1.5 border border-slate-700 uppercase tracking-wide cursor-pointer text-left font-mono {activeTab ===
      'logs'
        ? 'bg-slate-700 text-white'
        : 'bg-slate-800 text-slate-400 hover:text-slate-200'}"
    >
      [4] Event Logs
    </button>
  </nav>

  <!-- Main Content Area -->
  <main class="border border-slate-700 bg-slate-800 p-4 min-h-[340px] mb-4">
    {#if activeTab === 'nodes'}
      <div class="flex justify-between items-center mb-3">
        <h2 class="font-bold text-slate-200">Registered Manifest Nodes</h2>
        <button
          onclick={fetchNodes}
          class="text-xs bg-slate-700 hover:bg-slate-600 px-2 py-1 text-slate-200 cursor-pointer"
          >[Refresh]</button
        >
      </div>
      <table class="w-full text-left text-xs border border-slate-700">
        <thead>
          <tr class="bg-slate-900 border-b border-slate-700 text-slate-400">
            <th class="p-2 border-r border-slate-700">MANIFEST ID</th>
            <th class="p-2 border-r border-slate-700">NAME</th>
            <th class="p-2 border-r border-slate-700">CATEGORY</th>
            <th class="p-2 border-r border-slate-700">BINARY</th>
            <th class="p-2">VERSION</th>
          </tr>
        </thead>
        <tbody>
          {#each nodes as node}
            <tr class="border-b border-slate-700 hover:bg-slate-750">
              <td class="p-2 font-bold text-slate-300 border-r border-slate-700">{node.id}</td>
              <td class="p-2 border-r border-slate-700">{node.name}</td>
              <td class="p-2 border-r border-slate-700 text-slate-400">{node.category}</td>
              <td class="p-2 border-r border-slate-700 text-slate-300 font-bold"
                >{node.execution?.binary || 'internal'}</td
              >
              <td class="p-2 text-slate-400">{node.version}</td>
            </tr>
          {/each}
          {#if nodes.length === 0}
            <tr>
              <td colspan="5" class="p-4 text-center text-slate-500"
                >No node manifests loaded yet.</td
              >
            </tr>
          {/if}
        </tbody>
      </table>
    {/if}

    {#if activeTab === 'storage'}
      <div class="flex justify-between items-center mb-3">
        <h2 class="font-bold text-slate-200">Storage Pools & Pre-Flight Health Probes</h2>
        <button
          onclick={fetchStorage}
          class="text-xs bg-slate-700 hover:bg-slate-600 px-2 py-1 text-slate-200 cursor-pointer"
          >[Run Health Probes]</button
        >
      </div>
      <table class="w-full text-left text-xs border border-slate-700">
        <thead>
          <tr class="bg-slate-900 border-b border-slate-700 text-slate-400">
            <th class="p-2 border-r border-slate-700">ALIAS</th>
            <th class="p-2 border-r border-slate-700">HOST PATH</th>
            <th class="p-2 border-r border-slate-700">HEALTH STATUS</th>
            <th class="p-2">PROBE LATENCY</th>
          </tr>
        </thead>
        <tbody>
          {#each Object.entries(storagePools) as [alias, pool]}
            <tr class="border-b border-slate-700">
              <td class="p-2 font-bold border-r border-slate-700">{alias}</td>
              <td class="p-2 border-r border-slate-700 text-slate-300">{pool.path}</td>
              <td
                class="p-2 border-r border-slate-700 uppercase font-bold {pool.status === 'healthy'
                  ? 'text-slate-300'
                  : 'text-red-400'}"
              >
                [{pool.status}]
              </td>
              <td class="p-2 text-slate-400">{(pool.latency / 1000000).toFixed(2)} ms</td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}

    {#if activeTab === 'trigger'}
      <h2 class="font-bold text-slate-200 mb-3">Trigger File Job Execution</h2>
      <form onsubmit={handleTriggerJob} class="max-w-xl space-y-3 text-xs">
        <div>
          <label for="targetFile" class="block text-slate-400 mb-1">TARGET FILE PATH:</label>
          <input
            id="targetFile"
            type="text"
            bind:value={targetFilePath}
            class="w-full bg-slate-900 border border-slate-700 p-2 text-slate-100 font-mono focus:border-slate-500 focus:outline-none"
            placeholder="/path/to/target_file.mp4"
          />
        </div>
        <button
          type="submit"
          class="bg-slate-700 hover:bg-slate-600 text-slate-100 px-4 py-2 border border-slate-600 font-bold cursor-pointer"
        >
          EXECUTE FILE JOB
        </button>
        {#if triggerMessage}
          <div class="p-2 bg-slate-900 border border-slate-700 text-slate-300">
            {triggerMessage}
          </div>
        {/if}
      </form>
    {/if}

    {#if activeTab === 'logs'}
      <h2 class="font-bold text-slate-200 mb-2">Live Console Logs</h2>
      <div
        class="bg-slate-900 border border-slate-700 p-3 h-64 overflow-y-auto text-xs text-slate-300 space-y-1 font-mono"
      >
        {#each logLines as line}
          <div>{line}</div>
        {/each}
      </div>
    {/if}
  </main>

  <!-- Footer -->
  <footer class="text-xs text-slate-500 border-t border-slate-800 pt-2 flex justify-between">
    <span>CONDUIT SYSTEM // UTILITY DASHBOARD</span>
    <span>SVELTE 5 RUNES ENGAGED</span>
  </footer>
</div>
