export namespace service {
	
	export class Template {
	    // Go type: ulids.ULID
	    id: any;
	    name: string;
	    html: string;
	    // Go type: time.Time
	    createdAt: any;
	    // Go type: time.Time
	    updatedAt: any;
	    // Go type: gorm.DeletedAt
	    deletedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Template(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = this.convertValues(source["id"], null);
	        this.name = source["name"];
	        this.html = source["html"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.deletedAt = this.convertValues(source["deletedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Event {
	    // Go type: ulids.ULID
	    id: any;
	    // Go type: ulids.ULID
	    campaignID: any;
	    detail: string;
	    // Go type: time.Time
	    createdAt: any;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = this.convertValues(source["id"], null);
	        this.campaignID = this.convertValues(source["campaignID"], null);
	        this.detail = source["detail"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class File {
	    // Go type: ulids.ULID
	    id: any;
	    folder: string;
	    fileName: string;
	    // Go type: time.Time
	    createdAt: any;
	    // Go type: time.Time
	    updatedAt: any;
	    // Go type: gorm.DeletedAt
	    deletedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = this.convertValues(source["id"], null);
	        this.folder = source["folder"];
	        this.fileName = source["fileName"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.deletedAt = this.convertValues(source["deletedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Campaign {
	    // Go type: ulids.ULID
	    id: any;
	    // Go type: ulids.ULID
	    fileID?: any;
	    name: string;
	    body: string;
	    subject: string;
	    // Go type: ulids.ULID
	    templateID?: any;
	    // Go type: time.Time
	    createdAt: any;
	    // Go type: time.Time
	    updatedAt: any;
	    // Go type: gorm.DeletedAt
	    deletedAt: any;
	    file: File;
	    events: Event[];
	    template?: Template;
	
	    static createFrom(source: any = {}) {
	        return new Campaign(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = this.convertValues(source["id"], null);
	        this.fileID = this.convertValues(source["fileID"], null);
	        this.name = source["name"];
	        this.body = source["body"];
	        this.subject = source["subject"];
	        this.templateID = this.convertValues(source["templateID"], null);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.deletedAt = this.convertValues(source["deletedAt"], null);
	        this.file = this.convertValues(source["file"], File);
	        this.events = this.convertValues(source["events"], Event);
	        this.template = this.convertValues(source["template"], Template);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	

}

