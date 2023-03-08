let ct_p;
let t_p;
let mapCanvas;
let playersCanvas;
let homeElement, screenElement
const dimensions = 900
let imgFlash = new Image(), imgHe = new Image(), imgDecoy = new Image(), imgSmoke = new Image(), imgIncendiary = new Image(), imgMolo = new Image();

function loadNades() {
    imgFlash.src = "../img/nade_flash.webp";
    imgHe.src = "../img/nade_he.webp";
    imgDecoy.src = "../img/nade_decoy.webp";
    imgSmoke.src = "../img/nade_smoke.webp";
    imgMolo.src = "../img/nade_molo.webp";
    imgIncendiary.src = "../img/nade_incendiary.webp";
}

function init() {
    homeElement = document.getElementById("home");
    screenElement = document.getElementById("screen")
    mapCanvas = document.getElementById("game_map");
    mapCanvas.width = dimensions;
    mapCanvas.height = dimensions;
    document.getElementById("canvas_container").style.height = dimensions
    document.getElementById("canvas_container").style.width = dimensions;
    playersCanvas = document.getElementById("players_canvas");
    playersCanvas.width = dimensions;
    playersCanvas.height = dimensions;
    loadNades();
}

function drawMap(mapName) {
    homeElement.style.display = "none";
    screenElement.classList.remove('hidden');
    ctx_map = mapCanvas.getContext("2d");
    ctx_map.scale(dimensions / 1024, dimensions / 1024);
    const img = new Image();
    img.onload = () => {
        ctx_map.drawImage(img, 0, 0);
    };
    img.src = `../img/${mapName}_radar.png`;
};

function updateState(currentState) {
    ctx_pl = playersCanvas.getContext("2d");
    const state = JSON.parse(currentState);

    displayTimeLeft(state.RoundTimeLeft);

    teamT = state.TeamT;
    teamCt = state.TeamCt;

    document.getElementById("clan_tag_t").textContent = teamT.ClanTag;
    document.getElementById("clan_tag_ct").textContent = teamCt.ClanTag;
    document.getElementById("score_t").textContent = teamT.Score
    document.getElementById("score_ct").textContent = teamCt.Score

    ctx_pl.reset();
    ctx_pl.scale(dimensions / 1024, dimensions / 1024);
    for (let i = 0; i < teamT.Players.length + teamCt.Players.length; i++) {
        p = [...teamT.Players, ...teamCt.Players].sort(sortPlayersById)[i];
        if (!p.IsAlive) {
            drawDeadPlayer(ctx_pl, p);
        } else {
            drawAlivePlayer(ctx_pl, p);
        }
    }

    fillTable(teamT, teamCt);

    let nades = state.Nades;
    drawNades(ctx_pl, nades);
}

function drawNades(ctx, nades) {
    for (let i = 0; i < nades.length; i++) {
        let nade = nades[i];
        ctx.beginPath();
        ctx.moveTo(nade.Positions[0].X, nade.Positions[0].Y);
        let lastPosX, lastPosY;
        for (let j = 1; j < nade.Positions.length; j++) {
            pos = nade.Positions[j]
            ctx.lineTo(pos.X, pos.Y);
            ctx.moveTo(pos.X, pos.Y);
            lastPosX = pos.X;
            lastPosY = pos.Y;
        }
        ctx.strokeStyle = "white";
        ctx.closePath();
        ctx.stroke();
        nadeImg = getNadeImg(nade.Type);
        if (nadeImg !== null) {
            ctx.drawImage(nadeImg, lastPosX - 20, lastPosY - 20 / 2, 40, 40);
        }
    }
}

function getNadeImg(type) {
    if (type == "Decoy Grenade") {
        return imgDecoy;
    } else if (type == "Flashbang") {
        return imgFlash;
    } else if (type == "HE Grenade") {
        return imgHe;
    } else if (type == "Incendiary Grenade") {
        return imgIncendiary;
    } else if (type == "Molotov") {
        return imgMolo;
    } else if (type == "Smoke Grenade") {
        return imgSmoke;
    }

    return null
}

function sortPlayersByScore(a, b) {
    let scoreA = (a.Kills * 2) + a.Assists;
    let scoreB = (b.Kills * 2) + b.Assists;
    let diff = scoreB - scoreA;
    if (diff == 0) {
        return a.Id - b.Id
    }
    return diff
}

function sortPlayersById(a, b) {
    return a.Id - b.Id
}

function fillTable(teamT, teamCt) {
    teamT.Players.sort(sortPlayersByScore);
    for (let i = 0; i < teamT.Players.length; i++) {
        let row = document.getElementById(`p${i + 1}`)
        insertPlayerData(row, teamT.Players[i])
    }

    teamCt.Players.sort(sortPlayersByScore);
    for (let i = 0; i < teamCt.Players.length; i++) {
        let row = document.getElementById(`p${i + 6}`)
        insertPlayerData(row, teamCt.Players[i])
    }
}

function insertPlayerData(row, player) {
    row.querySelector("#name").textContent = `${player.Id} ${player.Name}`;
    row.querySelector("#hp").textContent = player.HP;
    row.querySelector("#money").textContent = player.Money;
    row.querySelector("#equipped").textContent = player.Equipped;
    row.querySelector("#kda").textContent = `${player.Kills}/${player.Deaths}/${player.Assists}`;
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
    let fillStyle;
    if (p.Team == 2) {
        fillStyle = "orange"
    } else {
        fillStyle = "#219ebc"
    }

    var hpDown = (100 - p.HP) / 100

    ctx.save();
    ctx.fillStyle = fillStyle;
    ctx.globalCompositeOperation = "destination-over";
    ctx.beginPath();
    ctx.translate(pos.X, pos.Y);
    ctx.rotate((- p.ViewDirection - 90) * Math.PI / 180.0);
    ctx.translate(-pos.X, -pos.Y);
    ctx.arc(pos.X, pos.Y, r, 0.5 * Math.PI + hpDown * Math.PI + 0.001, 0.5 * Math.PI - hpDown * Math.PI, false);
    ctx.closePath();
    ctx.fill();
    
    ctx.beginPath();
    ctx.fillStyle = "black";
    ctx.arc(pos.X, pos.Y, r, 0.5 * Math.PI + hpDown * Math.PI + 0.001, 0.5 * Math.PI - hpDown * Math.PI, true);
    
    ctx.closePath();
    ctx.fill();
    ctx.fillStyle = fillStyle;

    ctx.beginPath()
    ctx.strokeStyle = "white"
    ctx.lineWidth = 2
    ctx.arc(pos.X, pos.Y, r + 0.5, 0, 2 * Math.PI);
    ctx.stroke();
    ctx.closePath();

    ctx.moveTo(pos.X - r / 2, pos.Y + r + 2);
    ctx.lineTo(pos.X, pos.Y + r + r);
    ctx.lineTo(pos.X + r / 2, pos.Y + r + 2);
    ctx.fill();
    ctx.stroke();
    ctx.restore();
}