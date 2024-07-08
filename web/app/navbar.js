require("./index")

async function logout() {
    try {
        const response = await fetch("/auth", {
            method: "DELETE",
            // Set the FormData instance as the request body
            headers: {
                "Accept": "text/plain"
            }
        });
        if (response.status >= 400) {
            console.error("Unexpected error while logging out");
        } else if (response.status === 200) {
            window.location.href = "/auth";
        }
    } catch (e) {
        console.error(e);
    }
}

const logouts = document.querySelectorAll(".logout");
logouts.forEach(logoutEl => logoutEl.addEventListener("click", (event) => {
    event.preventDefault();
    logout();
}));