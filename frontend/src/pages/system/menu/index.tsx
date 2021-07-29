import { DeleteOutlined, PlusOutlined, FormOutlined } from '@ant-design/icons';
import { Button, Tooltip, Divider, Modal, message } from 'antd';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import React, { useState, useRef } from 'react';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { queryMenus, deleteMenu } from '@/services/anew/menu';
import { useAccess, Access } from 'umi';

const MenuList: React.FC = () => {
    const [createVisible, setCreateVisible] = useState<boolean>(false);
    const [updateVisible, setUpdateVisible] = useState<boolean>(false);
    const actionRef = useRef<ActionType>();
    const [formValues, setFormValues] = useState<API.MenuList>();
    const access = useAccess();

    const handleDelete = (record: API.Ids) => {
        if (!record) return;
        if (Array.isArray(record.ids) && !record.ids.length) return;
        const content = `您是否要删除这${Array.isArray(record.ids) ? record.ids.length : ''}项？`;
        Modal.confirm({
            title: '注意',
            content,
            onOk: () => {
                deleteMenu(record).then((res) => {
                    if (res.code === 200 && res.status === true) {
                        message.success(res.message);
                        if (actionRef.current) {
                            actionRef.current.reload();
                        }
                    }
                });
            },
            onCancel() { },
        });
    };

    const columns: ProColumns<API.MenuList>[] = [
        {
            title: '名称',
            dataIndex: 'name',
            search: false,
        },
        {
            title: '图标',
            dataIndex: 'icon',
            search: false,
        },
        {
            title: '路径',
            dataIndex: 'path',
            search: false,
        },
        {
            title: '排序',
            dataIndex: 'sort',
            search: false,
        },
        {
            title: '创建人',
            dataIndex: 'creator',
            search: false,
        },
        {
            title: '操作',
            dataIndex: 'option',
            valueType: 'option',
            render: (_, record: API.MenuList) => (
                <>
                    <Access accessible={access.hasPerms(['admin', 'menu:update'])}>
                        <Tooltip title="修改">
                            <FormOutlined
                                style={{ fontSize: '17px', color: '#52c41a' }}
                                onClick={() => {
                                    setFormValues(record);
                                    setUpdateVisible(true);
                                }}
                            />
                        </Tooltip>
                    </Access>
                    <Divider type="vertical" />
                    <Access accessible={access.hasPerms(['admin', 'menu:delete'])}>
                        <Tooltip title="删除">
                            <DeleteOutlined
                                style={{ fontSize: '17px', color: 'red' }}
                                onClick={() => handleDelete({ ids: [record.id] })}
                            />
                        </Tooltip>
                    </Access>
                </>
            ),
        },
    ];

    return (
        <PageHeaderWrapper>
            {/* 权限控制显示内容 */}
            {access.hasPerms(['admin', 'menu:list']) && <ProTable
                actionRef={actionRef}
                rowKey="id"
                pagination={false}
                search={false}
                toolBarRender={(action, { selectedRows }) => [
                    <Access accessible={access.hasPerms(['admin', 'menu:create'])}>
                        <Button key="1" type="primary" onClick={() => setCreateVisible(true)}>
                            <PlusOutlined /> 新建
                        </Button>
                    </Access>,
                    selectedRows && selectedRows.length > 0 && (
                        <Access accessible={access.hasPerms(['admin', 'menu:delete'])}>
                            <Button
                                key="2"
                                type="primary"
                                onClick={() => handleDelete({ ids: selectedRows.map((item) => item.id) })}
                                danger
                            >
                                <DeleteOutlined /> 删除
                            </Button>
                        </Access>
                    ),
                ]}
                tableAlertRender={({ selectedRowKeys, selectedRows }) => (
                    <div>
                        已选择{' '}
                        <a
                            style={{
                                fontWeight: 600,
                            }}
                        >
                            {selectedRowKeys.length}
                        </a>{' '}
                        项&nbsp;&nbsp;
                    </div>
                )}
                request={async (params) => queryMenus({ ...params })}
                columns={columns}
                rowSelection={{}}
            />}
            {createVisible && (
                <CreateForm
                    actionRef={actionRef}
                    onChange={setCreateVisible}
                    modalVisible={createVisible}
                />
            )}
            {updateVisible && (
                <UpdateForm
                    actionRef={actionRef}
                    onChange={setUpdateVisible}
                    modalVisible={updateVisible}
                    values={formValues}
                />
            )}
        </PageHeaderWrapper>
    );
};

export default MenuList;