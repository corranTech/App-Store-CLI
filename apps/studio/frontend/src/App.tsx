import { FormEvent, useState } from "react";

import "./styles.css";
import { ChatMessage, NavSection } from "./types";

const sections: NavSection[] = [
  { id: "overview", label: "Overview", description: "Release cockpit" },
  { id: "apps", label: "Apps", description: "App library" },
  { id: "builds", label: "Builds", description: "TestFlight and processing" },
  { id: "submission", label: "Submission", description: "Validation and publish" },
  { id: "assets", label: "Assets", description: "Screenshots and localizations" },
  { id: "threads", label: "Threads", description: "ACP history" },
];

const sectionIcons: Record<string, string> = {
  overview: "◎",
  apps: "⊞",
  builds: "⏣",
  submission: "↗",
  assets: "□",
  threads: "≡",
};

const initialMessages: ChatMessage[] = [];

const apps = [
  { name: "MusadoraKit", platform: "iOS", version: "2.3.0", status: "Ready for validation" },
  { name: "ASC Test", platform: "iOS", version: "1.7.2", status: "Safe sandbox app" },
  { name: "Composer Pad", platform: "macOS", version: "0.9.5", status: "Metadata drift" },
];

const releaseChecklist = [
  { label: "Build selected", status: "done", detail: "2.3.0 (451) from March 30" },
  { label: "App review details", status: "warning", detail: "Missing demo account notes" },
  { label: "Screenshots", status: "done", detail: "All 6.9″ and 5.5″ sets synced" },
  { label: "Localizations", status: "warning", detail: "French subtitle missing" },
  { label: "IAP readiness", status: "done", detail: "2 items approved" },
  { label: "Submission approval", status: "pending", detail: "Waiting on Studio confirmation" },
];

const buildRows = [
  { version: "2.3.0", build: "451", state: "Processing", age: "9m ago" },
  { version: "2.2.9", build: "447", state: "Ready for Sale", age: "3d ago" },
  { version: "2.2.8", build: "441", state: "Expired", age: "2w ago" },
];

const approvalPreview = [
  "asc validate --app 6759231657 --version 2.3.0 --platform IOS --output json",
  "asc publish appstore --app 6759231657 --version 2.3.0 --submit --confirm --output json",
];

const checkIcon: Record<string, { char: string; cls: string }> = {
  done: { char: "✓", cls: "check-done" },
  warning: { char: "⚠", cls: "check-warning" },
  pending: { char: "○", cls: "check-pending" },
};

export default function App() {
  const [activeSection, setActiveSection] = useState<NavSection>(sections[0]);
  const [messages, setMessages] = useState<ChatMessage[]>(initialMessages);
  const [draft, setDraft] = useState("");
  const [showLogs, setShowLogs] = useState(false);

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmed = draft.trim();
    if (!trimmed) {
      return;
    }

    setMessages((current) => [
      ...current,
      { id: `user-${current.length}`, role: "user", content: trimmed, timestamp: "Now" },
      {
        id: `assistant-${current.length}`,
        role: "assistant",
        content: "Bootstrap mode recorded the prompt. Next we can wire it into the live ACP transport instead of the placeholder response path.",
        timestamp: "Now",
      },
    ]);
    setDraft("");
  }

  return (
    <div className="studio-shell">
      {/* Sidebar */}
      <aside className="sidebar">
        <div className="sidebar-header">
          <span className="sidebar-title">ASC Studio</span>
          <button className="sidebar-action" type="button" aria-label="New thread">+</button>
        </div>

        <div className="sidebar-section">
          <p className="sidebar-section-label">Workspace</p>
          {sections.map((section) => (
            <button
              key={section.id}
              type="button"
              className={`sidebar-row ${section.id === activeSection.id ? "is-active" : ""}`}
              onClick={() => setActiveSection(section)}
            >
              <span className="sidebar-row-icon">{sectionIcons[section.id]}</span>
              <span>{section.label}</span>
            </button>
          ))}
        </div>

        <div className="sidebar-spacer" />

        <div className="thread-section">
          <p className="sidebar-section-label">Threads</p>
          <div className="thread-row is-selected">
            <strong>Ship 2.3.0</strong>
            <small>12m</small>
          </div>
          <div className="thread-row">
            <strong>Metadata sync</strong>
            <small>2d</small>
          </div>
        </div>
      </aside>

      <div className="shell-separator" />

      {/* Main area */}
      <div className="main-area">
        <header className="toolbar">
          <div className="toolbar-left">
            <span className="toolbar-breadcrumb">
              App Store Connect <span>/</span> {activeSection.label}
            </span>
          </div>
          <div className="toolbar-right">
            <span className="toolbar-status">
              <span className="status-dot" />
              Ready
            </span>
            <button className="toolbar-btn" type="button">Refresh</button>
          </div>
        </header>

        <section className="workspace">
          <div className="workspace-main">
            {/* Apps panel */}
            <div className="panel">
              <div className="panel-heading">
                <div>
                  <h3 className="panel-title">Apps</h3>
                  <p className="panel-subtitle">Release cockpit</p>
                </div>
              </div>
              <div className="app-list">
                {apps.map((app) => (
                  <button
                    key={app.name}
                    className={`app-row ${app.name === "MusadoraKit" ? "is-current" : ""}`}
                    type="button"
                  >
                    <div>
                      <strong>{app.name}</strong>
                      <small>{app.platform} {app.version}</small>
                    </div>
                    <span className="app-row-status">{app.status}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Builds panel */}
            <div className="panel">
              <div className="panel-heading">
                <div>
                  <h3 className="panel-title">Builds</h3>
                  <p className="panel-subtitle">Current release candidates</p>
                </div>
                <div className="segmented">
                  <button className="is-active" type="button">TestFlight</button>
                  <button type="button">App Store</button>
                </div>
              </div>
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Version</th>
                    <th>Build</th>
                    <th>State</th>
                    <th className="col-trailing">Age</th>
                  </tr>
                </thead>
                <tbody>
                  {buildRows.map((row) => (
                    <tr key={`${row.version}-${row.build}`}>
                      <td>{row.version}</td>
                      <td>{row.build}</td>
                      <td>
                        <span className={`build-state state-${row.state.toLowerCase().replace(/\s+/g, "-")}`}>
                          {row.state}
                        </span>
                      </td>
                      <td className="col-trailing">{row.age}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Inspector */}
          <div className="workspace-inspector">
            <div className="panel">
              <div className="panel-heading">
                <div>
                  <h3 className="panel-title">Submission readiness</h3>
                  <p className="panel-subtitle">2 blockers</p>
                </div>
              </div>
              <div className="checklist">
                {releaseChecklist.map((item) => {
                  const icon = checkIcon[item.status];
                  return (
                    <div key={item.label} className="checklist-row">
                      <span className={`check-icon ${icon.cls}`}>{icon.char}</span>
                      <div>
                        <strong>{item.label}</strong>
                        <small>{item.detail}</small>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>

            <div className="panel">
              <div className="panel-heading">
                <h3 className="panel-title">Mutation preview</h3>
              </div>
              <pre className="command-preview">{approvalPreview.join("\n")}</pre>
            </div>
          </div>
        </section>

        {/* Chat dock */}
        <section className="dock">
          {messages.length > 0 && (
            <div className="dock-body">
              <div className="message-list" aria-label="Chat messages">
                {messages.map((message) => (
                  <article key={message.id} className={`message-row role-${message.role}`}>
                    <p>{message.content}</p>
                  </article>
                ))}
              </div>

              {showLogs ? (
                <div className="logs-pane" data-testid="log-drawer">
                  <pre>
                    {`session/new {"cwd":"/Users/rudrank/Developer/CLIs/App-Store-Connect-CLI"}\n`}
                    {`session/prompt {"role":"user","content":"Validate version 2.3.0"}\n`}
                    {`session/update {"kind":"message","role":"assistant","content":"Validation in progress"}\n`}
                  </pre>
                </div>
              ) : null}
            </div>
          )}

          <form className="composer" onSubmit={handleSubmit}>
            <div className="composer-card">
              <textarea
                aria-label="Chat prompt"
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                placeholder="Ask Studio to inspect builds, explain blockers, or draft a command…"
                rows={2}
              />
              <div className="composer-bar">
                <div className="composer-meta">
                  <span>Codex</span>
                  <span>Claude</span>
                  <span>Custom ACP</span>
                </div>
                <button className="send-btn" type="submit" aria-label="Send">⬆</button>
              </div>
            </div>
          </form>
        </section>
      </div>
    </div>
  );
}
