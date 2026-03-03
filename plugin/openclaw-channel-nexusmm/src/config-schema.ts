import { z } from "zod";

const NexusmmAccountSchema = z.strictObject({
  name: z.string().optional(),
  enabled: z.boolean().optional(),
  botToken: z.string().optional(),
  apiUrl: z.string().optional(),
  pollIntervalMs: z.number().int().min(500).optional(),
});

export const NexusmmConfigSchema = z.strictObject({
  name: z.string().optional(),
  enabled: z.boolean().optional(),
  botToken: z.string().optional(),
  apiUrl: z.string().optional(),
  pollIntervalMs: z.number().int().min(500).optional(),
  accounts: z.record(z.string(), NexusmmAccountSchema.optional()).optional(),
});

export type NexusmmConfig = z.infer<typeof NexusmmConfigSchema>;
