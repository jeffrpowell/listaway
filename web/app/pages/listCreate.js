require('../navbar')

document.addEventListener('DOMContentLoaded', (event) => {
    const listNameInputs = document.getElementsByName('name');
    const submitButtons = document.querySelectorAll('button');
    const errorSpans = document.querySelectorAll('.error-span');
    var formReadyToSubmit = false;

    listNameInputs.forEach(listNameInput => 
        listNameInput.addEventListener('input', (event) => {
            submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
            errorSpans.forEach(errorSpan => {
                errorSpan.classList.add('hidden');
                errorSpan.textContent = '';
            });
            if (listNameInput.value.trim() !== '') {
                debouncedCheckListName(submitButtons, errorSpans, listNameInput.value.trim());
            }
            else {
                formReadyToSubmit = false;
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
            }
        })
    );

    var debouncedCheckListName = debounce(checkListName, 500);
    async function checkListName(submitButtons, errorSpans, name) {
        try {
            const response = await fetch(`/list/namecheck?name=${encodeURIComponent(name)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.status === 200) {
                submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
                formReadyToSubmit = true;
            } else if (response.status === 400) {
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'The list name is already in use.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            } else {
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            }
        } catch (error) {
            submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
            errorSpans.forEach(errorSpan => {
                errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                errorSpan.classList.remove('hidden');
            });
            formReadyToSubmit = false;
        }
    }
    async function sendData(form) {
        if (!formReadyToSubmit) {
            return;
        }
        // Associate the FormData object with the form element
        const formData = new FormData(form);

        try {
            const response = await fetch("/list", {
                method: "PUT",
                headers: {
                    "Accept": "text/plain",
                    "Content-Type": "application/x-www-form-urlencoded"
                },
                body: new URLSearchParams(formData).toString()
            });
            if (response.status >= 400) {
                showError(response.status);
            } else if (response.status === 200) {
                window.location.href = response.headers.get("Location");
            }
        } catch (e) {
            console.error(e);
            showError(500);
        }
    }

    // Take over form submission
    const forms = document.querySelectorAll(".list-form");
    forms.forEach(form => form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendData(form);
    }));

    function showError(statusCode) {
        const errorSpans = document.querySelectorAll(".error-span");
        errorSpans.forEach(errorSpan => {
            if (statusCode === 400) {
                errorSpan.innerText = "List name is already in use.";
            }
            else {
                errorSpan.innerText = "Unexpected error occurred. Please try again later.";
            }
            errorSpan.classList.remove("hidden");
        });
    }

    function debounce(func, delay) {
        let timeoutId;
        return function(...args) {
            if (timeoutId) {
                clearTimeout(timeoutId);
            }
            timeoutId = setTimeout(() => {
                func.apply(this, args);
            }, delay);
        };
    }
});