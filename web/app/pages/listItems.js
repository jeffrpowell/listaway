require('../index')
require('../navbar')
import { themeQuartz } from '@ag-grid-community/theming';
import { ClientSideRowModelModule } from "@ag-grid-community/client-side-row-model";
import { createGrid } from "@ag-grid-community/core";

const itemNameRenderer = function(params) {
    if (params.data.url.Valid) {
        return `<a href="${params.data.url.String}" class="underline" target="_blank" rel="nofollow noreferrer">${params.value}</a>`
    } else {
        return params.value;
    }
}

const priorityComparator = function(p1, p2) {
    const p1Empty = p1 === undefined || p1 === null || p1 === "";
    const p2Empty = p2 === undefined || p2 === null || p2 === "";
    if (p1Empty && p2Empty) {
        return 0;
    } else if (p1Empty) {
        return 1; //place p1 after p2
    } else if (p2Empty) {
        return -1; //place p2 after p1
    } else {
        return p1 - p2;
    }
}

class ActionButtonsComponent {
    eGui;
    deleteButton;
    eventListener;

    init(params) {
        this.eGui = document.createElement('div');
        this.eGui.classList.add('flex');
        this.eGui.innerHTML = 
            "<a href='/list/"+params.context.listId+"/item/"+params.data.id+"/edit'>"
            // + "<!-- https://heroicons.com/ pencil -->"
            +     "<svg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke-width='1.5' stroke='currentColor' class='size-6 text-gray-600'>"
            +         "<path stroke-linecap='round' stroke-linejoin='round' d='m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125' />"
            +     "</svg>"
            + "</a>";
        this.deleteButton = document.createElement('button');
        this.deleteButton.setAttribute("type", "button");
        this.deleteButton.dataset.listId = params.context.listId;
        this.deleteButton.dataset.itemId = params.data.id;
        this.deleteButton.innerHTML = 
            // +  <!-- https://heroicons.com/ trash -->"
            "<svg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke-width='1.5' stroke='currentColor' class='size-6 text-red-900'>"
            +     "<path stroke-linecap='round' stroke-linejoin='round' d='m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0' />"
            + "</svg>"
        this.eventListener = async event => {
            let listId = this.deleteButton.dataset.listId;
            let itemId = this.deleteButton.dataset.itemId;
            const response = await fetch('/list/' + listId + '/item/' + itemId, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
            });
            if (response.status === 204 || response.status === 200) {
                window.location.href = '/list/' + listId;
            }
        }
        this.deleteButton.addEventListener('click', this.eventListener);
        this.eGui.appendChild(this.deleteButton);
    }

    getGui() {
        return this.eGui;
    }

    refresh() {
        return true;
    }

    destroy() {
        if (this.deleteButton) {
            this.deleteButton.removeEventListener('click', this.eventListener);
        }
    }
}

document.addEventListener('DOMContentLoaded', (event) => {
    const grid = document.querySelectorAll('.item-grid');
    const customTheme = themeQuartz.withParams({
        // internalContentLineHeight: "1.5"
    });

    const gridOptions = {
        theme: customTheme,
        loadThemeGoogleFonts: true,
        rowData: [],
        columnDefs: [
            { 
                field: "name", 
                headerName: "Item",
                flex: 1,
                minWidth: 120,
                wrapText: true,
                autoHeight: true,
                cellRenderer: itemNameRenderer
            },
            { 
                field: "priority", 
                headerName: "Priority", 
                type: 'numericColumn',
                width: 115,
                sort: 'asc',
                valueGetter: p => p.data.priority.Valid ? p.data.priority.Int64 : null,
                comparator: priorityComparator
            },
            { 
                field: "notes", 
                headerName: "Notes", 
                flex: 2,
                minWidth: 240,
                wrapText: true,
                autoHeight: true,
                valueGetter: p => p.data.notes.Valid ? p.data.notes.String : ""
            },
            {
                field: "actions",
                headerName: "Actions",
                width: 100,
                cellRenderer: ActionButtonsComponent
            }
        ],
        domLayout: 'autoHeight',
        getRowId: params => String(params.data.id),
        autoSizeStrategy: {
            type: "fitGridWidth",
        }
    };
    grid.forEach(g => {
        let listId = g.dataset.listId;
        let augmentedGridOptions = gridOptions;
        augmentedGridOptions.context = {listId: listId};
        const gridApi = createGrid(g, augmentedGridOptions, {
            modules: [
                ClientSideRowModelModule,
            ]
        });
        fetch("/list/"+listId+"/items")
            .then(response => response.json())
            .then(data => gridApi.setGridOption("rowData", data));
    });
});