import { fireEvent, render, screen } from "@testing-library/react";

import App from "./App";

describe("App", () => {
  it("switches workspace sections from the sidebar", () => {
    render(<App />);

    fireEvent.click(screen.getByRole("button", { name: /builds/i }));

    expect(screen.getByRole("heading", { name: "Builds" })).toBeInTheDocument();
    expect(screen.getByText(/current release candidates/i)).toBeInTheDocument();
  });

  it("sends a chat message and shows response", () => {
    render(<App />);

    const textarea = screen.getByLabelText("Chat prompt");
    fireEvent.change(textarea, { target: { value: "list builds" } });
    fireEvent.submit(textarea.closest("form")!);

    expect(screen.getByText("list builds")).toBeInTheDocument();
    expect(screen.getByText(/bootstrap mode/i)).toBeInTheDocument();
  });
});
