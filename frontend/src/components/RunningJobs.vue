<template>
  <div>
    <div class="flex items-center justify-between mb-4 px-2">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold">Running Jobs</h2>
        <span class="badge badge-success badge-soft badge-sm">Live</span>
      </div>
      <div class="flex items-center gap-3">
        <span v-if="lastUpdated" class="text-xs text-base-content/50">
          Updated {{ lastUpdatedLabel }}
        </span>
        <button
          class="btn btn-sm btn-ghost"
          :disabled="loading"
          @click="fetchJobs"
        >
          <span
            v-if="loading"
            class="loading loading-spinner loading-xs"
          ></span>
          <span v-else>↻ Refresh</span>
        </button>
      </div>
    </div>

    <div v-if="loading && jobs.length === 0" class="flex justify-center py-12">
      <span
        class="loading loading-spinner loading-md text-base-content/30"
      ></span>
    </div>

    <div
      v-else-if="jobs.length === 0"
      class="text-center py-12 text-base-content/40"
    >
      No jobs currently running.
    </div>

    <div
      v-else
      class="overflow-x-auto rounded-box border border-base-content/10"
    >
      <table class="table table-zebra w-full">
        <thead>
          <tr>
            <th>Name</th>
            <th>Run ID</th>
            <th>Started</th>
            <th>Runtime</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="job in jobs" :key="job.id">
            <td class="font-medium">{{ job.name }}</td>
            <td>
              <span :title="job.id" class="font-mono text-sm">
                {{ job.id.slice(0, 8) }}
              </span>
            </td>
            <td class="text-sm">{{ formatDate(job.started_at) }}</td>
            <td class="font-mono text-sm tabular-nums">
              {{ runtimes[job.id] ?? "00:00:00" }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { toUtcDate, formatDate, formatDuration } from "../utils/time";

interface Job {
  id: string;
  name: string;
  started_at: string;
}

const jobs = ref<Job[]>([]);
const loading = ref(false);
const lastUpdated = ref<Date | null>(null);
const runtimes = ref<Record<string, string>>({});

let pollTimer: ReturnType<typeof setInterval> | null = null;
let tickTimer: ReturnType<typeof setInterval> | null = null;

async function fetchJobs() {
  loading.value = true;
  try {
    const res = await fetch("/web/jobs/running");
    if (res.ok) {
      jobs.value = await res.json();
      lastUpdated.value = new Date();
    }
  } finally {
    loading.value = false;
  }
}

function formatElapsed(startedAt: string): string {
  const elapsed = Math.max(
    0,
    Math.floor((Date.now() - toUtcDate(startedAt).getTime()) / 1000),
  );
  return formatDuration(elapsed);
}

const lastUpdatedLabel = ref("");
function updateLastUpdatedLabel() {
  if (!lastUpdated.value) return;
  const secs = Math.floor((Date.now() - lastUpdated.value.getTime()) / 1000);
  if (secs < 10) lastUpdatedLabel.value = "just now";
  else if (secs < 60) lastUpdatedLabel.value = `${secs}s ago`;
  else lastUpdatedLabel.value = `${Math.floor(secs / 60)}m ago`;
}

function tick() {
  for (const job of jobs.value) {
    runtimes.value[job.id] = formatElapsed(job.started_at);
  }
  updateLastUpdatedLabel();
}

onMounted(async () => {
  await fetchJobs();
  tick();
  tickTimer = setInterval(tick, 1000);
  pollTimer = setInterval(fetchJobs, 30_000);
});

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer);
  if (tickTimer) clearInterval(tickTimer);
});
</script>
