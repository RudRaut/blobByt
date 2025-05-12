module sui_vault::sui_vault {

    use sui::table::{ Table };
    use sui::table;
    use sui::clock::Clock;


    public struct File has key, store {
        id: UID,
        uploader: address,
        blob_id: address,
        timestamp: u64,
        // access_list: Table<address, bool>,
    }

    // public struct DownloadLog has store {

    // }

    // public struct Signature has store {

    // }

    public struct FileRegistry has key, store {
        id: UID,                                      // Unique ID for the file registry
        files_by_user: Table<address, vector<File>>, // Files uploaded or accessible by user
        // download_log: Table<ID, vector<DownloadLog>>, // FileID => list of download events
        // signatures: Table<ID, vector<Signature>>, // FileID => list of signatures
    }

    fun init(ctx: &mut TxContext) {
        let registry =FileRegistry {
            id: object::new(ctx),
            files_by_user: table::new<address, vector<File>>(ctx),
            // download_log: table::new<ID, vector<DownloadLog>>(ctx),
            // signatures: table::new<ID, vector<Signature>>(ctx)
        };
        transfer::transfer(registry, tx_context::sender(ctx))
    }

    public entry fun upload_file(registry: &mut FileRegistry, blob_id: address, clock: &Clock, ctx: &mut TxContext) {
        let uploader = tx_context::sender(ctx);
        let file_id = object::new(ctx);
        let timestamp = clock.timestamp_ms();

        let file = File {
            id: file_id,
            uploader: uploader,
            blob_id,
            timestamp: timestamp,
        };
        if (!table::contains(&registry.files_by_user, uploader)){
            table::add(&mut registry.files_by_user, uploader, vector::empty<File>());
        };

        let user_files = table::borrow_mut(&mut registry.files_by_user, uploader);
        vector::push_back(user_files, file);
    }

    
}

