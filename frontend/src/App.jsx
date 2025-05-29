import React, { useState } from "react";

const Sidebar = ({ selected, setSelected }) => (
    <div style={{
        width: 140,
        background: "#222",
        color: "#fff",
        height: "100vh",
        paddingTop: 40,
        position: "fixed",
        left: 0,
        top: 0,
        display: "flex",
        flexDirection: "column",
        gap: 10
    }}>
        {['create', 'update', 'delete', 'stats'].map((item) => (
            <button
                key={item}
                style={{
                    background: selected === item ? "#007bff" : "#333",
                    color: "#fff",
                    border: "none",
                    borderRadius: 4,
                    padding: "12px 0",
                    margin: "0 10px",
                    fontWeight: "bold",
                    cursor: "pointer"
                }}
                onClick={() => setSelected(item)}
            >
                {item.charAt(0).toUpperCase() + item.slice(1)}
            </button>
        ))}
    </div>
);

const CreateForm = ({ setResult }) => {
    const [url, setUrl] = useState("");
    const [code, setCode] = useState("");
    return (
        <form
            onSubmit={async (e) => {
                e.preventDefault();
                setResult("");
                const body = code ? { url, short_code: code } : { url };
                try {
                    const res = await fetch("/shorten", {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(body),
                    });
                    const data = await res.json();
                    if (res.ok) {
                        setResult(
                            <div className="result">Short URL: <a href={data.short_url} target="_blank" rel="noopener noreferrer">{data.short_url}</a></div>
                        );
                    } else {
                        setResult(<div className="error">{data.error || `Error: ${res.status}`}</div>);
                    }
                } catch {
                    setResult(<div className="error">Network error. Please try again.</div>);
                }
            }}
            style={{ display: "flex", flexDirection: "column", gap: 12, maxWidth: 400 }}
        >
            <label>Long URL</label>
            <input type="url" value={url} onChange={e => setUrl(e.target.value)} placeholder="https://example.com" required />
            <label>Custom Short Code (optional)</label>
            <input type="text" value={code} onChange={e => setCode(e.target.value)} placeholder="e.g. mycode" />
            <button type="submit">Shorten</button>
        </form>
    );
};

const UpdateForm = ({ setResult }) => {
    const [code, setCode] = useState("");
    const [url, setUrl] = useState("");
    const [newCode, setNewCode] = useState("");
    return (
        <form
            onSubmit={async (e) => {
                e.preventDefault();
                setResult("");
                const body = newCode ? { url, short_code: newCode } : { url };
                try {
                    const res = await fetch(`/u/${code}`, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(body),
                    });
                    if (res.ok) {
                        setResult(<div className="result">Short URL updated successfully.</div>);
                    } else {
                        const data = await res.json().catch(() => ({}));
                        setResult(<div className="error">{data.error || `Error: ${res.status}`}</div>);
                    }
                } catch {
                    setResult(<div className="error">Network error. Please try again.</div>);
                }
            }}
            style={{ display: "flex", flexDirection: "column", gap: 12, maxWidth: 400 }}
        >
            <label>Short Code to Update</label>
            <input type="text" value={code} onChange={e => setCode(e.target.value)} placeholder="e.g. exmp" required />
            <label>New Long URL</label>
            <input type="url" value={url} onChange={e => setUrl(e.target.value)} placeholder="https://new-url.com" required />
            <label>New Short Code (optional)</label>
            <input type="text" value={newCode} onChange={e => setNewCode(e.target.value)} placeholder="e.g. newcode" />
            <button type="submit">Update</button>
        </form>
    );
};

const DeleteForm = ({ setResult }) => {
    const [code, setCode] = useState("");
    return (
        <form
            onSubmit={async (e) => {
                e.preventDefault();
                setResult("");
                try {
                    const res = await fetch(`/u/${code}`, { method: "DELETE" });
                    if (res.ok) {
                        setResult(<div className="result">Short URL deleted successfully.</div>);
                    } else {
                        setResult(<div className="error">Delete failed. Error: {res.status}</div>);
                    }
                } catch {
                    setResult(<div className="error">Network error. Please try again.</div>);
                }
            }}
            style={{ display: "flex", flexDirection: "column", gap: 12, maxWidth: 400 }}
        >
            <label>Short Code to Delete</label>
            <input type="text" value={code} onChange={e => setCode(e.target.value)} placeholder="e.g. exmp" required />
            <button type="submit" style={{ background: "#dc3545", color: "#fff" }}>Delete</button>
        </form>
    );
};

const StatsForm = ({ setResult }) => {
    const [code, setCode] = useState("");
    return (
        <form
            onSubmit={async (e) => {
                e.preventDefault();
                setResult("");
                try {
                    const res = await fetch(`/stats/${code}`);
                    if (res.ok) {
                        const data = await res.json();
                        setResult(<div className="result">Access count: {data.access_count}</div>);
                    } else {
                        setResult(<div className="error">Stats not found. Error: {res.status}</div>);
                    }
                } catch {
                    setResult(<div className="error">Network error. Please try again.</div>);
                }
            }}
            style={{ display: "flex", flexDirection: "column", gap: 12, maxWidth: 400 }}
        >
            <label>Short Code for Stats</label>
            <input type="text" value={code} onChange={e => setCode(e.target.value)} placeholder="e.g. exmp" required />
            <button type="submit">Get Stats</button>
        </form>
    );
};

const Main = () => {
    const [selected, setSelected] = useState("create");
    const [result, setResult] = useState("");
    let content;
    if (selected === "create") content = <CreateForm setResult={setResult} />;
    if (selected === "update") content = <UpdateForm setResult={setResult} />;
    if (selected === "delete") content = <DeleteForm setResult={setResult} />;
    if (selected === "stats") content = <StatsForm setResult={setResult} />;
    return (
        <div style={{ display: "flex" }}>
            <Sidebar selected={selected} setSelected={setSelected} />
            <div style={{ marginLeft: 160, padding: 40, width: "100%" }}>
                <h1 style={{ color: "#333" }}>URL Shortener</h1>
                {content}
                <div style={{ marginTop: 24 }}>{result}</div>
            </div>
        </div>
    );
};

export default Main;
