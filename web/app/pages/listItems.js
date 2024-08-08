require('../index')
require('../navbar')

const deleteItemButtons = document.querySelectorAll('.btn-delete-item');

deleteItemButtons.forEach(deleteItem => {
    deleteItem.addEventListener('click', async (event) => {
        let listId = deleteItem.dataset.listId;
        let itemId = deleteItem.dataset.itemId;
        const response = await fetch('/list/'+listId+'/item/'+itemId, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
        });
        if (response.status === 204 || response.status === 200) {
            window.location.href = '/list/'+listId;
        }
    });
});