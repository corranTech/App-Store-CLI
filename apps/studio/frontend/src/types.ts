export type NavSection = {
  id: string;
  label: string;
  description: string;
};

export type ChatMessage = {
  id: string;
  role: "user" | "assistant" | "system";
  content: string;
  timestamp: string;
};

export type ApprovalCard = {
  id: string;
  title: string;
  summary: string;
  commandPreview: string[];
  status: "pending" | "approved" | "rejected";
};
