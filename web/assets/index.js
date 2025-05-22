if (localStorage.getItem("theme") === '2') {
    document.body.classList.add('cyberpunk');
} else {
    document.body.classList.remove('cyberpunk');
}

fetch('/json', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            engine: localStorage.getItem("engine")
        })
    })
    .then(res => res.text())
    .then(console.log)
    .catch(console.error);


const selectTheme = document.getElementById("themes");
selectTheme.addEventListener('change', function() {
    const theme = selectTheme.value;
    if (theme === '2') {
        document.body.classList.add('cyberpunk');
    } else {
        document.body.classList.remove('cyberpunk');
    }
    localStorage.setItem("theme", selectTheme.value);
})

const selectEngine = document.getElementById("engines")
selectEngine.addEventListener('change', function() {
    localStorage.setItem("engine", selectEngine.value);
})
