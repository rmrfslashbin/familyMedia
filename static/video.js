(function () {
    if (Hls.isSupported()) {
        var video = document.getElementById('video');
        var hls = new Hls();
        // bind them together
        hls.attachMedia(video);
        hls.on(Hls.Events.MEDIA_ATTACHED, function () {
            const videoElement = document.querySelector('#video');
            hls.loadSource(videoElement.dataset.src);
            hls.on(Hls.Events.MANIFEST_PARSED, function (event, data) {
                video.play()
            });
        });
        hls.on(Hls.Events.ERROR, function (event, data) {
            if (data.fatal) {
                switch (data.type) {
                    case Hls.ErrorTypes.NETWORK_ERROR:
                        // try to recover network error
                        console.log('fatal network error encountered, try to recover');
                        hls.startLoad();
                        break;
                    case Hls.ErrorTypes.MEDIA_ERROR:
                        console.log('fatal media error encountered, try to recover');
                        hls.recoverMediaError();
                        break;
                    default:
                        // cannot recover
                        hls.destroy();
                        break;
                }
            }
        });
    }
})();