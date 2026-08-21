package api

import (
	"net/http"
)

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Unborn — Population</title>
  <style>
    :root { --bg:#0f1115; --card:#1a1d24; --text:#e8eaed; --muted:#9aa0a6; --accent:#7c9cff; --good:#3dd68c; --warn:#f5a524; --bad:#f76e6e; }
    * { box-sizing: border-box; }
    body { margin:0; font-family: ui-sans-serif, system-ui, sans-serif; background:var(--bg); color:var(--text); }
    header { padding:1.25rem 1.5rem; border-bottom:1px solid #2a2f3a; display:flex; align-items:center; justify-content:space-between; gap:1rem; flex-wrap:wrap; }
    header h1 { margin:0; font-size:1.25rem; font-weight:600; letter-spacing:-0.02em; }
    header span { color:var(--muted); font-size:0.875rem; }
    main { padding:1.5rem; max-width:1100px; margin:0 auto; }
    .strip { display:grid; grid-template-columns:repeat(auto-fit,minmax(140px,1fr)); gap:0.75rem; margin-bottom:1.5rem; }
    .stat { background:var(--card); border-radius:10px; padding:1rem; }
    .stat .label { color:var(--muted); font-size:0.75rem; text-transform:uppercase; letter-spacing:0.04em; }
    .stat .value { font-size:1.5rem; font-weight:600; margin-top:0.25rem; }
    table { width:100%; border-collapse:collapse; background:var(--card); border-radius:10px; overflow:hidden; }
    th, td { text-align:left; padding:0.75rem 1rem; border-bottom:1px solid #2a2f3a; font-size:0.9rem; }
    th { color:var(--muted); font-weight:500; font-size:0.75rem; text-transform:uppercase; letter-spacing:0.04em; }
    tr:last-child td { border-bottom:none; }
    .pill { display:inline-block; padding:0.15rem 0.5rem; border-radius:999px; font-size:0.75rem; font-weight:500; }
    .thriving { background:#1a3d2e; color:var(--good); }
    .stable { background:#1e2a3d; color:var(--accent); }
    .under_pressure { background:#3d2e1a; color:var(--warn); }
    .critical, .collapsed { background:#3d1a1a; color:var(--bad); }
    .actions { margin-bottom:1.5rem; display:flex; gap:0.5rem; flex-wrap:wrap; }
    button { background:var(--accent); color:#0f1115; border:none; border-radius:8px; padding:0.5rem 0.9rem; font-weight:600; cursor:pointer; font-size:0.875rem; }
    button.secondary { background:#2a2f3a; color:var(--text); }
    button:disabled { opacity:0.5; cursor:not-allowed; }
    .empty { color:var(--muted); padding:2rem; text-align:center; }
    code { font-size:0.8rem; color:var(--muted); }
  </style>
</head>
<body>
  <header>
    <div>
      <h1>Unborn</h1>
      <span>Population overview · farming with souls</span>
    </div>
    <span id="clock"></span>
  </header>
  <main>
    <div class="strip" id="strip">
      <div class="stat"><div class="label">Personas</div><div class="value" id="s-personas">—</div></div>
      <div class="stat"><div class="label">Bodies</div><div class="value" id="s-bodies">—</div></div>
      <div class="stat"><div class="label">Thriving</div><div class="value" id="s-thriving">—</div></div>
      <div class="stat"><div class="label">Pressure+</div><div class="value" id="s-pressure">—</div></div>
    </div>
    <div class="actions">
      <button type="button" id="btn-refresh">Refresh</button>
      <button type="button" class="secondary" id="btn-create">Create sample persona</button>
    </div>
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Location</th>
          <th>Engagement</th>
          <th>Vitality</th>
          <th>Level</th>
          <th>ID</th>
        </tr>
      </thead>
      <tbody id="rows">
        <tr><td colspan="6" class="empty">Loading…</td></tr>
      </tbody>
    </table>
  </main>
  <script>
    const $ = (id) => document.getElementById(id);
    function levelClass(level) {
      return level || 'stable';
    }
    async function load() {
      const [personas, vitality, instances] = await Promise.all([
        fetch('/v1/personas').then(r => r.json()),
        fetch('/v1/vitality').then(r => r.json()),
        fetch('/v1/instances').then(r => r.json()),
      ]);
      const vmap = {};
      (vitality || []).forEach(v => { vmap[v.persona_id] = v; });
      let thriving = 0, pressure = 0;
      (vitality || []).forEach(v => {
        const lv = v.value >= 80 ? 'thriving' : v.value >= 55 ? 'stable' : v.value >= 30 ? 'under_pressure' : v.value >= 10 ? 'critical' : 'collapsed';
        if (lv === 'thriving') thriving++;
        if (['under_pressure','critical','collapsed'].includes(lv)) pressure++;
      });
      $('s-personas').textContent = (personas || []).length;
      $('s-bodies').textContent = (instances || []).length;
      $('s-thriving').textContent = thriving;
      $('s-pressure').textContent = pressure;
      const tbody = $('rows');
      if (!personas || !personas.length) {
        tbody.innerHTML = '<tr><td colspan="6" class="empty">No personas yet. Create one to begin.</td></tr>';
        return;
      }
      tbody.innerHTML = personas.map(p => {
        const v = vmap[p.id];
        const val = v ? v.value.toFixed(1) : '75.0';
        const lv = v ? (v.value >= 80 ? 'thriving' : v.value >= 55 ? 'stable' : v.value >= 30 ? 'under_pressure' : v.value >= 10 ? 'critical' : 'collapsed') : 'stable';
        return '<tr>' +
          '<td>' + (p.display_name || '—') + '</td>' +
          '<td>' + (p.demographics && p.demographics.location || '—') + '</td>' +
          '<td>' + (p.engagement && p.engagement.type || '—') + '</td>' +
          '<td>' + val + '</td>' +
          '<td><span class="pill ' + levelClass(lv) + '">' + lv + '</span></td>' +
          '<td><code>' + p.id.slice(0,8) + '…</code></td>' +
          '</tr>';
      }).join('');
    }
    $('btn-refresh').onclick = () => load();
    $('btn-create').onclick = async () => {
      $('btn-create').disabled = true;
      await fetch('/v1/personas', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          display_name: 'Persona ' + Math.floor(Math.random()*1000),
          location: 'Berlin',
          timezone: 'Europe/Berlin',
          age_min: 25, age_max: 30,
          engagement: 'thoughtful_commenter'
        })
      });
      $('btn-create').disabled = false;
      load();
    };
    setInterval(() => { $('clock').textContent = new Date().toLocaleString(); }, 1000);
    load();
    setInterval(load, 15000);
  </script>
</body>
</html>`

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}
