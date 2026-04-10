/** Robyn serializes naive datetimes without a timezone suffix — append Z so
 *  the browser treats them as UTC rather than local time. */
export function toUtcDate(iso: string): Date {
  return new Date(iso.endsWith("Z") ? iso : iso + "Z");
}

export function formatDate(iso: string): string {
  return toUtcDate(iso).toLocaleString();
}

export function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  return `${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
}
