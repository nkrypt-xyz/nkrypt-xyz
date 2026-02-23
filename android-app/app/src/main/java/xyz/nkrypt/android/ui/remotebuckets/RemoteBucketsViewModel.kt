package xyz.nkrypt.android.ui.remotebuckets

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import xyz.nkrypt.android.data.local.entity.RemoteBucketEntity
import xyz.nkrypt.android.data.remote.RemoteBucketRepository
import javax.inject.Inject

@HiltViewModel
class RemoteBucketsViewModel @Inject constructor(
    private val repository: RemoteBucketRepository
) : ViewModel() {

    private val _buckets = MutableStateFlow<List<RemoteBucketEntity>>(emptyList())
    val buckets: StateFlow<List<RemoteBucketEntity>> = _buckets.asStateFlow()

    init {
        viewModelScope.launch {
            repository.getAllBuckets().collect { list ->
                _buckets.value = list
            }
        }
    }

    fun removeBucket(id: String) {
        viewModelScope.launch {
            repository.deleteBucket(id)
        }
    }
}
