require("./index")

async function logout() {
    const response = await fetch("/auth", {
        method: "DELETE",
        headers: {
            "Accept": "text/plain"
        },
    });
    if (response.status === 204 || response.status === 200) {
        window.location.href = response.headers.get("Location");
    }  
    else if (response.status >= 400) {
        console.error("Unexpected error while logging out");
    }
}

const logouts = document.querySelectorAll(".logout");
logouts.forEach(logoutEl => logoutEl.addEventListener("click", (event) => {
    logout();
}));