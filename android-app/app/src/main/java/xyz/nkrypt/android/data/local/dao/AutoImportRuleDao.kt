package xyz.nkrypt.android.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow
import xyz.nkrypt.android.data.local.entity.AutoImportRuleEntity

@Dao
interface AutoImportRuleDao {

    @Query("SELECT * FROM auto_import_rules ORDER BY createdAt DESC")
    fun getAll(): Flow<List<AutoImportRuleEntity>>

    @Query("SELECT * FROM auto_import_rules WHERE id = :id")
    suspend fun getById(id: String): AutoImportRuleEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entity: AutoImportRuleEntity)

    @Query("DELETE FROM auto_import_rules WHERE id = :id")
    suspend fun deleteById(id: String)
}
