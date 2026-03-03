import type { PluginRuntime } from "openclaw/plugin-sdk";

let runtime: PluginRuntime | null = null;

export function setNexusmmRuntime(next: PluginRuntime) {
  runtime = next;
}

export function getNexusmmRuntime(): PluginRuntime {
  if (!runtime) {
    throw new Error("Nexus-MM runtime not initialized");
  }
  return runtime;
}
