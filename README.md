## IDEA 
Have a bugeting software that sits as an independet app and almost never connects to the internet i.e , have the app UI, server and DB sitted on the same device.

## TODO 
- Enforce params for filter endpoints

# Critical TODO 
- Switch from mutable params to structs

# To uderstand why categories was modeled as such , check this
https://web.archive.org/web/20180729174436/http://www.tomjewett.com/dbdesign/dbdesign.php?page=recursive.php



Why not do uuid.UUID `json:"association_id" gorm:"type:uuid;default:uuid_generate_v4()"` ?
To avoid this having to run direct DB queries 
See issue below : 

The extension is available but not installed in this database.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";