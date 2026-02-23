package xyz.nkrypt.android.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import xyz.nkrypt.android.data.local.entity.LocalBlobEntity

@Dao
interface LocalBlobDao {

    @Query("SELECT * FROM local_blobs WHERE fileId = :fileId ORDER BY createdAt DESC LIMIT 1")
    suspend fun getLatestByFileId(fileId: String): LocalBlobEntity?

    @Query("SELECT * FROM local_blobs WHERE id = :id")
    suspend fun getById(id: String): LocalBlobEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entity: LocalBlobEntity)

    @Query("DELETE FROM local_blobs WHERE fileId = :fileId")
    suspend fun deleteByFileId(fileId: String)
}
