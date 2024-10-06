import 'w2ui/w2ui-2.0.css';
import { w2grid } from 'w2ui/w2ui-2.0.es6';
require('../index')
require('../navbar')

document.addEventListener('DOMContentLoaded', (event) => {
    const deleteItemButtons = document.querySelectorAll('.btn-delete-item');
    const w2uigrid = document.querySelectorAll('.w2grid');
    
    const grid = new w2grid({
        name: 'listItemsGrid',
        fixedBody: false,
        columns: [
            { field: 'item', text: 'Item', size: '100%' },
            { field: 'delete', text: 'Delete', size: '50px' }
        ],
        records: [
            // TODO: spin up new endpoint to feed this data to the table
            { recid: 1, item: 'Sample Item 1', delete: '<button class="btn-delete-item" data-list-id="1" data-item-id="1">Delete</button>' },
            { recid: 2, item: 'Sample Item 2', delete: '<button class="btn-delete-item" data-list-id="1" data-item-id="2">Delete</button>' },
        ]
    });
    w2uigrid.forEach(g => grid.render(g));
    
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
});