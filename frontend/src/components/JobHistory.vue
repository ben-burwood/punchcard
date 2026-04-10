<template>
  <div>
    <div class="flex items-center justify-between mb-4 px-2">
      <h2 class="text-lg font-semibold">History</h2>
      <input
        v-model="search"
        type="search"
        placeholder="Search by name…"
        class="input input-sm input-bordered w-56"
      />
    </div>

    <div v-if="loading && items.length === 0" class="flex justify-center py-12">
      <span
        class="loading loading-spinner loading-md text-base-content/30"
      ></span>
    </div>

    <div
      v-else-if="!loading && items.length === 0"
      class="text-center py-12 text-base-content/40"
    >
      No completed jobs found.
    </div>

    <template v-else>
      <div class="overflow-x-auto rounded-box border border-base-content/10">
        <table class="table table-zebra w-full">
          <thead>
            <tr>
              <th>
                <button
                  class="btn btn-ghost btn-xs gap-1"
                  @click="setSort('name')"
                >
                  Name
                  <span class="text-xs opacity-40">{{ sortIcon("name") }}</span>
                </button>
              </th>
              <th>Run ID</th>
              <th>
                <button
                  class="btn btn-ghost btn-xs gap-1"
                  @click="setSort('started_at')"
                >
                  Started
                  <span class="text-xs opacity-40">{{
                    sortIcon("started_at")
                  }}</span>
                </button>
              </th>
              <th>
                <button
                  class="btn btn-ghost btn-xs gap-1"
                  @click="setSort('stopped_at')"
                >
                  Stopped
                  <span class="text-xs opacity-40">{{
                    sortIcon("stopped_at")
                  }}</span>
                </button>
              </th>
              <th>
                <button
                  class="btn btn-ghost btn-xs gap-1"
                  @click="setSort('duration')"
                >
                  Duration
                  <span class="text-xs opacity-40">{{
                    sortIcon("duration")
                  }}</span>
                </button>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="job in items" :key="job.id">
              <td class="font-medium">{{ job.name }}</td>
              <td>
                <span :title="job.id" class="font-mono text-sm">{{
                  job.id.slice(0, 8)
                }}</span>
              </td>
              <td class="text-sm">{{ formatDate(job.started_at) }}</td>
              <td class="text-sm">{{ formatDate(job.stopped_at) }}</td>
              <td class="font-mono text-sm tabular-nums">
                {{ formatDuration(job.duration_seconds) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex items-center justify-between mt-3 px-1">
        <span class="text-sm text-base-content/50">
          {{ total }} result{{ total === 1 ? "" : "s" }}
        </span>
        <div class="flex items-center gap-2">
          <button
            class="btn btn-sm btn-ghost"
            :disabled="page === 1"
            @click="page--"
          >
            ← Prev
          </button>
          <span class="text-sm">Page {{ page }} of {{ totalPages }}</span>
          <button
            class="btn btn-sm btn-ghost"
            :disabled="page >= totalPages"
            @click="page++"
          >
            Next →
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { toUtcDate, formatDate, formatDuration } from "../utils/time";

interface HistoryJob {
  id: string;
  name: string;
  started_at: string;
  stopped_at: string;
  duration_seconds: number;
}

const PAGE_SIZE = 25;

const items = ref<HistoryJob[]>([]);
const total = ref(0);
const loading = ref(false);
const search = ref("");
const sort = ref("stopped_at");
const order = ref<"asc" | "desc">("desc");
const page = ref(1);

const totalPages = computed(() =>
  Math.max(1, Math.ceil(total.value / PAGE_SIZE)),
);

let debounceTimer: ReturnType<typeof setTimeout> | null = null;

async function fetchHistory() {
  loading.value = true;
  try {
    const params = new URLSearchParams({
      search: search.value,
      sort: sort.value,
      order: order.value,
      limit: String(PAGE_SIZE),
      offset: String((page.value - 1) * PAGE_SIZE),
    });
    const res = await fetch(`/web/jobs/history?${params}`);
    if (res.ok) {
      const data = await res.json();
      items.value = data.items;
      total.value = data.total;
    }
  } finally {
    loading.value = false;
  }
}

function setSort(col: string) {
  if (sort.value === col) {
    order.value = order.value === "asc" ? "desc" : "asc";
  } else {
    sort.value = col;
    order.value = "desc";
  }
  page.value = 1;
}

function sortIcon(col: string): string {
  if (sort.value !== col) return "↕";
  return order.value === "asc" ? "↑" : "↓";
}

watch(search, () => {
  if (debounceTimer) clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    page.value = 1;
    fetchHistory();
  }, 300);
});

watch([sort, order, page], fetchHistory);

fetchHistory();
</script>
