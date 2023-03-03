var canvasElement;
var ctx;

function init() {
    drawMap()
}

function drawMap() {
    canvasElement = document.getElementById("game_map");
    canvasElement.width = window.innerHeight
    canvasElement.height = window.innerHeight

    ctx = canvasElement.getContext("2d");

    const img = new Image();
    img.onload = () => {
        // context.canvas.style.imageRendering = 'auto';
        ctx.drawImage(img, 0, 0);
    };
    img.src = "../img/de_ancient_radar.png";
};

function drawPlayer(player) {
    console.log(player)
}

