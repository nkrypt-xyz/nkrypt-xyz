package xyz.nkrypt.android.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow
import xyz.nkrypt.android.data.local.entity.RemoteBucketEntity

@Dao
interface RemoteBucketDao {

    @Query("SELECT * FROM remote_buckets ORDER BY createdAt DESC")
    fun getAll(): Flow<List<RemoteBucketEntity>>

    @Query("SELECT * FROM remote_buckets WHERE id = :id")
    suspend fun getById(id: String): RemoteBucketEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entity: RemoteBucketEntity)

    @Query("UPDATE remote_buckets SET cachedApiKey = :apiKey, apiKeyExpiresAt = :expiresAt WHERE id = :id")
    suspend fun updateCachedApiKey(id: String, apiKey: String, expiresAt: Long?)

    @Query("DELETE FROM remote_buckets WHERE id = :id")
    suspend fun deleteById(id: String)
}
