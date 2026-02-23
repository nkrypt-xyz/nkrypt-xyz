package xyz.nkrypt.android.data.local

import androidx.room.Database
import androidx.room.RoomDatabase
import xyz.nkrypt.android.data.local.dao.AutoImportRuleDao
import xyz.nkrypt.android.data.local.dao.AutoSyncRuleDao
import xyz.nkrypt.android.data.local.dao.LocalBlobDao
import xyz.nkrypt.android.data.local.dao.LocalBucketDao
import xyz.nkrypt.android.data.local.dao.LocalDirectoryDao
import xyz.nkrypt.android.data.local.dao.LocalFileDao
import xyz.nkrypt.android.data.local.dao.RemoteBucketDao
import xyz.nkrypt.android.data.local.entity.AutoImportRuleEntity
import xyz.nkrypt.android.data.local.entity.AutoSyncRuleEntity
import xyz.nkrypt.android.data.local.entity.LocalBlobEntity
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity
import xyz.nkrypt.android.data.local.entity.LocalDirectoryEntity
import xyz.nkrypt.android.data.local.entity.LocalFileEntity
import xyz.nkrypt.android.data.local.entity.RemoteBucketEntity

@Database(
    entities = [
        LocalBucketEntity::class,
        LocalDirectoryEntity::class,
        LocalFileEntity::class,
        LocalBlobEntity::class,
        RemoteBucketEntity::class,
        AutoImportRuleEntity::class,
        AutoSyncRuleEntity::class
    ],
    version = 5,
    exportSchema = false
)
abstract class AppDatabase : RoomDatabase() {
    abstract fun localBucketDao(): LocalBucketDao
    abstract fun localDirectoryDao(): LocalDirectoryDao
    abstract fun localFileDao(): LocalFileDao
    abstract fun localBlobDao(): LocalBlobDao
    abstract fun remoteBucketDao(): RemoteBucketDao
    abstract fun autoImportRuleDao(): AutoImportRuleDao
    abstract fun autoSyncRuleDao(): AutoSyncRuleDao
}
