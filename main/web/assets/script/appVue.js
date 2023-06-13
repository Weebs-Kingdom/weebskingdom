let main = new Vue({
    'el': '#app',
    data: {
        loaded: false,
        late: false
    },
    created: async function () {
        this.initApp();
        // check if content loaded and show another button
        setTimeout(() => {
            if (!this.loaded) {
                this.late = true;
            }
        }, 1000 * 5)
    },
    methods: {
        reload() {
            window.location.reload()
        },
        home() {
            window.location.replace("/")
        },
        initApp: async function () {
            var fetchUri = document.location.pathname + "app";
            if (document.location.pathname == "/") {
                fetchUri = "/.app";
            }
            await fetch(fetchUri, {
                method: 'GET',
            }).then(async response => {
                document.getElementsByTagName("weeb-app")[0].innerHTML = await response.text();
                initColorModeBtn();
                this.loaded = true;
            })
        }
    }
})