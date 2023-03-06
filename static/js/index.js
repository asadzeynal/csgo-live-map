let ct_p;
let t_p;
let mapCanvas;
let playersCanvas;
let homeElement, screenElement
const dimensions = 900

function init() {
    homeElement = document.getElementById("home")
    screenElement = document.getElementById("screen")
    mapCanvas = document.getElementById("game_map");
    mapCanvas.width = dimensions
    mapCanvas.height = dimensions
    document.getElementById("canvas_container").style.height = dimensions
    document.getElementById("canvas_container").style.width = dimensions
    playersCanvas = document.getElementById("players_canvas");
    playersCanvas.width = dimensions
    playersCanvas.height = dimensions
}

function drawMap(mapName) {
    homeElement.style.display = "none"
    screenElement.classList.remove('hidden')
    ctx_map = mapCanvas.getContext("2d");
    ctx_map.scale(dimensions / 1024, dimensions / 1024)
    const img = new Image();
    img.onload = () => {
        ctx_map.drawImage(img, 0, 0);
    };
    img.src = `../img/${mapName}_radar.png`;
};

function updateState(currentState) {
    ctx_pl = playersCanvas.getContext("2d");
    const state = JSON.parse(currentState)

    displayTimeLeft(state.RoundTimeLeft)

    teamT = state.TeamT
    teamCt = state.TeamCt

    document.getElementById("clan_tag_t").textContent = teamT.ClanTag
    document.getElementById("clan_tag_ct").textContent = teamCt.ClanTag

    ctx_pl.reset()
    ctx_pl.scale(dimensions / 1024, dimensions / 1024)
    for (let i = 0; i < teamT.Players.length + teamCt.Players.length; i++) {
        p = [...teamT.Players, ...teamCt.Players][i]
        if (p.IsAlive) {
            drawAlivePlayer(ctx_pl, p)
        } else {
            drawDeadPlayer(ctx_pl, p)
        }
    }

    fillTable(teamT, teamCt)

    let nades = state.Nades
    drawNades(ctx_pl, nades)
}

function drawNades(ctx, nades) {
    for (let i = 0; i < nades.length; i++) {
        let nade = nades[i]
        ctx.beginPath();
        ctx.moveTo(nade.Positions[0].X, nade.Positions[0].Y)
        for (let j = 1; j < nade.Positions.length; j++) {
            pos = nade.Positions[j]
            ctx.lineTo(pos.X, pos.Y);
            ctx.moveTo(pos.X, pos.Y);
        }
        ctx.strokeStyle = "white"
        ctx.closePath()
        ctx.stroke()
    }
}

function sortPlayers(a, b) {
    let scoreA = (a.Kills * 2) + a.Assists
    let scoreB = (b.Kills * 2) + b.Assists
    let diff = scoreB - scoreA;
    if (diff == 0) {
        return a.Id - b.Id
    }
    return diff
}

function fillTable(teamT, teamCt) {
    teamT.Players.sort(sortPlayers);
    for (let i = 0; i < teamT.Players.length; i++) {
        let p = teamT.Players[i]
        let row = document.getElementById(`p${i + 1}`)
        row.querySelector("#name").textContent = `${p.Id} ${p.Name}`
        row.querySelector("#money").textContent = p.Money
        row.querySelector("#equipped").textContent = p.Equipped
        row.querySelector("#kda").textContent = `${p.Kills}/${p.Deaths}/${p.Assists}`
    }

    teamCt.Players.sort(sortPlayers);
    for (let i = 0; i < teamCt.Players.length; i++) {
        let p = teamCt.Players[i]
        let row = document.getElementById(`p${i + 6}`)
        row.querySelector("#name").textContent = `${p.Id} ${p.Name}`
        row.querySelector("#money").textContent = p.Money
        row.querySelector("#equipped").textContent = p.Equipped
        row.querySelector("#kda").textContent = `${p.Kills}/${p.Deaths}/${p.Assists}`
    }
}

function displayTimeLeft(duration) {
    let timeElement = document.getElementById("time_left")
    let minutes = Math.floor(duration / 60)
    let seconds = duration % 60
    if (seconds < 10) {
        seconds = `0${seconds}`
    }
    timeElement.textContent = `${minutes}:${seconds}`
}

function drawDeadPlayer(ctx, p) {
    pos = p.LastAlivePosition
    if (p.Team == 2) {
        ctx.strokeStyle = "orange"
    } else {
        ctx.strokeStyle = "#219ebc"
    }
    ctx.beginPath();
    ctx.moveTo(pos.X - 5, pos.Y - 5);
    ctx.lineTo(pos.X + 5, pos.Y + 5);
    ctx.moveTo(pos.X + 5, pos.Y - 5);
    ctx.lineTo(pos.X - 5, pos.Y + 5);
    ctx.closePath()
    ctx.stroke()
}

function drawAlivePlayer(ctx, p) {
    const r = 7
    pos = p.Position
    if (p.Team == 2) {
        ctx.fillStyle = "orange"
    } else {
        ctx.fillStyle = "#219ebc"
    }

    ctx.save();
    ctx.translate(pos.X, pos.Y)
    ctx.rotate((- p.ViewDirection - 90) * Math.PI / 180.0);
    ctx.translate(-pos.X, -pos.Y)
    ctx.beginPath();
    ctx.arc(pos.X, pos.Y, r, 0, 2 * Math.PI);
    ctx.moveTo(pos.X - r / 2, pos.Y + r + 2);
    ctx.lineTo(pos.X, pos.Y + r + r);
    ctx.lineTo(pos.X + r / 2, pos.Y + r + 2);
    ctx.closePath()
    ctx.fill();
    ctx.restore();
}