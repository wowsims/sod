type TableSorterRowData = {
	readonly values: ReadonlyArray<string | number>;
	readonly rowElement: HTMLTableRowElement;
};

type TableSorterConfig = {
	tableHead: HTMLTableRowElement;
	tableBody: HTMLTableSectionElement;
	dataSetKey: string;
	childRowClass: string;
	defaultSortCol: number;
	defaultSortDesc: boolean;
};

export class TableSorter {
	private readonly cfg: Readonly<TableSorterConfig>;
	private readonly rowData: Array<TableSorterRowData & { children?: Array<TableSorterRowData> }> = [];
	private sortCol = -1;
	private sortDesc: Array<boolean>;

	constructor(config: TableSorterConfig) {
        this.cfg = config;

		this.sortCol = this.cfg.defaultSortCol;
		this.sortDesc = Array(config.tableHead.cells.length).fill(true);
        this.sortDesc[config.defaultSortCol] = config.defaultSortDesc;

		for (let i = 0; i < config.tableHead.cells.length; i++) {
			config.tableHead.cells[i].addEventListener('click', () => this.setSort(i));
		}
	}

	private sortFunc = (a: TableSorterRowData, b: TableSorterRowData) => {
		const aValue = a.values[this.sortCol];
		const bValue = b.values[this.sortCol];
		const asc = !this.sortDesc[this.sortCol];
		if (typeof aValue === 'number' && typeof bValue === 'number') {
			return asc ? aValue - bValue : bValue - aValue;
		} else {
			return asc ? aValue.toString().localeCompare(bValue.toString()) : bValue.toString().localeCompare(aValue.toString());
		}
	};

	private sort() {
		if (this.rowData.length === 0) return;
		if (typeof this.rowData[0].values[this.sortCol] === 'undefined') return;

		this.rowData.sort(this.sortFunc);
		for (const row of this.rowData) {
			if (row.children) row.children.sort(this.sortFunc);
		}

        const body = this.cfg.tableBody;
		body.innerHTML = '';
		for (const row of this.rowData) {
			body.appendChild(row.rowElement);
			if (row.children) {
				for (const childRow of row.children) {
					body.appendChild(childRow.rowElement);
				}
			}
		}
	}

	/**
	 * Set column to sort by. If set to the current sort column the order will be reversed.
	 * @param column If omitted use default column.
	 */
	setSort(column?: number) {
		if (typeof column !== 'number') column = this.cfg.defaultSortCol;
		this.sortDesc[column] = !this.sortDesc[column];
		this.sortCol = column;
		this.sort();
	}

	private parseRowValues(rowElement: HTMLTableRowElement): Array<number | string> {
		const values: Array<string | number> = [];
		for (const cell of rowElement.cells) {
			const val = cell.dataset[this.cfg.dataSetKey] ?? cell.innerText;
			const numVal = parseFloat(val);
			values.push(!isNaN(numVal) ? numVal : val);
		}
		return values;
	}

	/**
	 * Update internal data structure for changed table data.
	 */
	update() {
		this.rowData.length = 0;

		for (const rowElement of this.cfg.tableBody.rows) {
			const values = this.parseRowValues(rowElement);
			if (!rowElement.classList.contains(this.cfg.childRowClass)) {
				this.rowData.push({ values, rowElement });
			} else {
				const parentData = this.rowData[this.rowData.length - 1];
				if (!parentData) throw new Error("Child row has no parent!");
				if (!parentData.children) parentData.children = [];
				parentData.children.push({ values, rowElement });
			}
		}

		this.sort();
	}
}
