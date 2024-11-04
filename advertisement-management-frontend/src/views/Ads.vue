<template>
	<div>
		<h1>广告管理</h1>
		<button @click="showAddAdForm = true">添加广告</button>

		<!-- 广告列表 -->
		<table>
			<thead>
				<tr>
					<th>标题</th>
					<th>描述</th>
					<th>图片</th>
					<th>目标链接</th>
					<th>状态</th>
					<th>操作</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="ad in ads" :key="ad.id">
					<td>{{ ad.title }}</td>
					<td>{{ ad.description }}</td>
					<td><img :src="ad.image_url" alt="广告图片" width="100" /></td>
					<td>
						<a :href="ad.target_url" target="_blank">{{ ad.target_url }}</a>
					</td>
					<td>{{ ad.status }}</td>
					<td>
						<button @click="editAd(ad)">编辑</button>
						<button @click="deleteAd(ad.id)">删除</button>
					</td>
				</tr>
			</tbody>
		</table>

		<!-- 添加/编辑广告表单 -->
		<div v-if="showAddAdForm || showEditAdForm">
			<h2>{{ showAddAdForm ? '添加广告' : '编辑广告' }}</h2>
			<form @submit.prevent="submitAdForm">
				<div>
					<label>标题：</label>
					<input v-model="form.title" required />
				</div>
				<div>
					<label>描述：</label>
					<textarea v-model="form.description" required></textarea>
				</div>
				<div>
					<label>图片URL：</label>
					<input v-model="form.image_url" required />
				</div>
				<div>
					<label>目标链接：</label>
					<input v-model="form.target_url" required />
				</div>
				<div>
					<label>状态：</label>
					<select v-model="form.status" required>
						<option value="active">激活</option>
						<option value="inactive">非激活</option>
					</select>
				</div>
				<button type="submit">{{ showAddAdForm ? '添加' : '更新' }}</button>
				<button type="button" @click="resetForm">取消</button>
			</form>
		</div>
	</div>
</template>

<script>
	import api from '@/services/api';

	export default {
		data() {
			return {
				ads: [],
				showAddAdForm: false,
				showEditAdForm: false,
				form: {
					id: null,
					title: '',
					description: '',
					image_url: '',
					target_url: '',
					status: 'active',
				},
			};
		},
		methods: {
			async fetchAds() {
				try {
					const response = await api.get('/ads');
					this.ads = response.data;
				} catch (error) {
					console.error('获取广告列表失败:', error);
				}
			},
			async submitAdForm() {
				try {
					if (this.showAddAdForm) {
						await api.post('/ads', this.form);
					} else if (this.showEditAdForm) {
						await api.put(`/ads/${this.form.id}`, this.form);
					}
					this.resetForm();
					this.fetchAds();
				} catch (error) {
					console.error('提交表单失败:', error);
				}
			},
			editAd(ad) {
				this.form = { ...ad };
				this.showEditAdForm = true;
			},
			async deleteAd(id) {
				try {
					await api.delete(`/ads/${id}`);
					this.fetchAds();
				} catch (error) {
					console.error('删除广告失败:', error);
				}
			},
			resetForm() {
				this.form = {
					id: null,
					title: '',
					description: '',
					image_url: '',
					target_url: '',
					status: 'active',
				};
				this.showAddAdForm = false;
				this.showEditAdForm = false;
			},
		},
		created() {
			this.fetchAds();
		},
	};
</script>

<style scoped>
	table {
		width: 100%;
		border-collapse: collapse;
		margin-top: 20px;
	}
	th,
	td {
		border: 1px solid #ddd;
		padding: 8px;
	}
	th {
		background-color: #f2f2f2;
	}
	form div {
		margin-bottom: 10px;
	}
</style>
