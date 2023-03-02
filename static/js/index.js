var canvasElement;
var ctx;

function init() {
    drawMap()
    initFileSelector()
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

function initFileSelector() {
    const fileSelector = document.getElementById('file-selector');
    fileSelector.addEventListener('change', (event) => {
        const fileList = event.target.files;
        setDemoFile(fileList[0])
    });
}
