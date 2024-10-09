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
        show: {
            toolbar: true,
            toolbarAdd: true,
            toolbarEdit: true,
            toolbarDelete: true,
            selectColumn: true
        },
        columns: [
            { field: 'item', text: 'Item', sortable: true, resizeable: false, searchable: false},
            { field: 'url', text: 'URL', sortable: true, resizeable: false, searchable: false},
            { field: 'priority', text: 'Priority', sortable: true, resizeable: false, searchable: false},
            { field: 'notes', text: 'Notes', sortable: true, resizeable: false, searchable: false}
        ],
        records: [
            // TODO: spin up new endpoint to feed this data to the table
            { recid: 1, item: 'Sample Item 1', url: '', priority: '', notes: '', w2ui: {class: {item: 'font-sans text-base', url: 'font-sans text-base', priority: 'font-sans text-base', notes: 'font-sans text-base'}}},
            { recid: 2, item: 'Sample Item 2', url: '', priority: '', notes: '', w2ui: {class: {item: 'font-sans text-base', url: 'font-sans text-base', priority: 'font-sans text-base', notes: 'font-sans text-base'}}},
            { recid: 3, item: 'Sample Item 3', url: 'https://jeffpowell.dev', priority: '', notes: '', w2ui: {class: {item: 'font-sans text-base', url: 'font-sans text-base', priority: 'font-sans text-base', notes: 'font-sans text-base'}}  },
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