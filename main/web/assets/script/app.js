async function initApp(){
    var fetchUri = document.location.pathname + "app";
    if (document.location.pathname == "/") {
        fetchUri = "/.app";
    }
    await fetch(fetchUri, {
        method: 'GET',
    }).then(async response => {
        document.getElementsByTagName("weeb-app")[0].innerHTML = await response.text();
    })
}