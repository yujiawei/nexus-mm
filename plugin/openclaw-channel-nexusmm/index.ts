/**
 * openclaw-channel-nexusmm
 *
 * OpenClaw channel plugin for Nexus-MM messaging platform.
 * Connects via HTTP polling (getUpdates) for real-time messaging.
 */

import type { OpenClawPluginApi } from "openclaw/plugin-sdk";
import { nexusmmPlugin } from "./src/channel.js";
import { setNexusmmRuntime } from "./src/runtime.js";

const plugin: {
  id: string;
  name: string;
  description: string;
  register: (api: OpenClawPluginApi) => void;
} = {
  id: "nexusmm",
  name: "Nexus-MM",
  description: "OpenClaw Nexus-MM channel plugin via HTTP polling",
  register(api) {
    setNexusmmRuntime(api.runtime);
    api.registerChannel({ plugin: nexusmmPlugin });
  },
};

export default plugin;
