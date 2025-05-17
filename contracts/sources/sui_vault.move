module sui_vault::sui_vault {

    use sui::table::{ Table };
    use sui::table;
    use sui::clock::Clock;

    public struct File has store {
        id: ID,
        uploader: address,
        blob_id: address,
        timestamp: u64,
        // access list needs to be implemented in FileRegistry
    }

    public struct DownloadLog has store {
        user: address,
        timestamp: u64
    }

    // public struct Signature has store {
    //     signer: address,
    //     signature: vector<u8>,
    //     timestamp: u64,
    // }

    public struct FileRegistry has key {
        id: UID,                                      // Unique ID for the file registry
        files: Table<ID, File>,
        files_by_user: Table<address, vector<ID>>, // Files uploaded or accessible by user
        access_lists: Table<ID, Table<address, bool>>, // Nested Table, FileID => address (True -> access granted)
        download_log: Table<ID, vector<DownloadLog>>, // FileID => list of download events
        // signatures: Table<ID, vector<Signature>>, // FileID => list of signatures
    }

    fun init(ctx: &mut TxContext) {
        let registry = FileRegistry {
            id: object::new(ctx),
            files: table::new<ID, File>(ctx),
            files_by_user: table::new<address, vector<ID>>(ctx),
            access_lists: table::new<ID, Table<address, bool>>(ctx),
            download_log: table::new<ID, vector<DownloadLog>>(ctx),
            // signatures: table::new<ID, vector<Signature>>(ctx)
        };
        transfer::transfer(registry, tx_context::sender(ctx))
    }

    // File Upload
    public entry fun upload_file(registry: &mut FileRegistry, blob_id: address, clock: &Clock, ctx: &mut TxContext) {
        let uploader = tx_context::sender(ctx);
        let file_id_addr = tx_context::fresh_object_address(ctx);
        let file_id =  object::id_from_address(file_id_addr);
        let timestamp = clock.timestamp_ms();
        let access_list = table::new<address, bool>(ctx);
        table::add(&mut registry.access_lists, file_id, access_list);

        let file = File {
            id: file_id,
            uploader: uploader,
            blob_id,
            timestamp: timestamp,
        };

        table::add(&mut registry.files, file_id, file);
        if (!table::contains(&registry.files_by_user, uploader)){
            table::add(&mut registry.files_by_user, uploader, vector::empty<ID>());
        };

        let user_files = table::borrow_mut(&mut registry.files_by_user, uploader);
        vector::push_back(user_files, file_id);

        grant_access(registry, file_id, uploader);
        has_access(registry, file_id, uploader);
    }


    public fun has_access(
        registry: &FileRegistry,
        file_id: ID,
        user: address
    ): bool {
        if (table::contains(&registry.access_lists, file_id)) {
            let access_list = table::borrow(&registry.access_lists, file_id);
            if (table::contains(access_list, user)) {
                let access = table::borrow(access_list, user);
                return *access
            };
        };
        false
    }


    public fun grant_access(registry: &mut FileRegistry, file_id: ID, user: address) {
        if (!has_access(registry, file_id, user)) {
            let access_list = table::borrow_mut(&mut registry.access_lists, file_id);
            if(!table::contains(access_list, user)) {
                table::add(access_list, user, true);
            }
        }
    }

    public fun revoke_access(registry: &mut FileRegistry, file_id: ID, user: address) {
        if(!has_access(registry, file_id, user)) {
            let access_list = table::borrow_mut(&mut registry.access_lists, file_id);
            if(table::contains(access_list, user)) {
                table::remove(access_list, user);
            }
        }
    }

    public fun list_user_files(registry: &FileRegistry, user: address): (vector<ID>, vector<ID>) {
        let mut owned = vector::empty<ID>();
        let mut accessible = vector::empty<ID>();
        
        if (table::contains(&registry.files_by_user, user)) {
            let user_files = table::borrow(&registry.files_by_user, user);
            let len = vector::length(user_files);
            let mut i = 0;

            while (i < len) {
                let file_id = *vector::borrow(user_files, i);
                let file = table::borrow(&registry.files, file_id);

                if(file.uploader == user) {
                    vector::push_back(&mut owned, file_id);
                } else if (table::contains(&registry.access_lists, file_id)) {
                    let access_table = table::borrow(&registry.access_lists, file_id);
                    if (table::contains(access_table, user)) {
                        let has_access = *table::borrow(access_table, user);
                        if (has_access) {
                            vector::push_back(&mut accessible, file_id);
                        }
                    }
                };

                i = i + 1;
            }
        };
        (owned, accessible)
    }

    public fun log_download(
        registry: &mut FileRegistry,
        file_id: ID,
        clock: &Clock,
        ctx: &mut TxContext
    ) {
        let user = tx_context::sender(ctx);
        let timestamp = clock.timestamp_ms();

        let log_entry = DownloadLog {
            user,
            timestamp,
        };

        if (table::contains(&registry.download_log, file_id)) {
            let logs = table::borrow_mut(&mut registry.download_log, file_id);
            vector::push_back(logs, log_entry);
        } else {
            let logs = vector::singleton(log_entry);
            table::add(&mut registry.download_log, file_id, logs);
        }
    }

}