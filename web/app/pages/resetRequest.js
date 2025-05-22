document.addEventListener('DOMContentLoaded', function () {
    const form = document.querySelector('.reset-request-form');
    if (!form) return;
    form.addEventListener('submit', async function (e) {
        e.preventDefault();
        const email = form.querySelector('input[name="email"]').value;
        const errorSpan = form.querySelector('.reset-error-span');
        errorSpan.classList.add('hidden');
        errorSpan.textContent = '';
        try {
            const res = await fetch('/reset', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ email }),
            });
            if (res.ok) {
                form.innerHTML = '<p class="text-green-600">If your email is registered, a reset link has been sent.</p>';
            } else {
                const msg = await res.text();
                errorSpan.textContent = msg || 'Failed to send reset link.';
                errorSpan.classList.remove('hidden');
            }
        } catch (err) {
            errorSpan.textContent = 'Network error.';
            errorSpan.classList.remove('hidden');
        }
    });
});
