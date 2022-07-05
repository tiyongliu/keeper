export namespace bridge {
	
	export class RefreshRequest {
	    conid: string;
	    keepOpen: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RefreshRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.conid = source["conid"];
	        this.keepOpen = source["keepOpen"];
	    }
	}

}

