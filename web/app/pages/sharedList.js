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
        ],
        domLayout: 'autoHeight',
        getRowId: params => String(params.data.id),
        autoSizeStrategy: {
            type: "fitGridWidth",
        }
    };
    grid.forEach(g => {
        let shareCode = g.dataset.shareCode;
        const gridApi = createGrid(g, gridOptions, {
            modules: [
                ClientSideRowModelModule,
            ]
        });
        fetch("/sharedlist/"+shareCode+"/items")
            .then(response => response.json())
            .then(data => gridApi.setGridOption("rowData", data));
    });
});