import React, { useEffect, useState } from 'react';
import { updatePermsRole } from '../service';
import { queryMenus } from '@/services/anew/menu';
import { queryApis } from '@/services/anew/apis';
import { getRolePermsByID } from '@/services/anew/role';
import { DrawerForm } from '@ant-design/pro-form';
import { message, Tree, Checkbox, Col, Row, Divider } from 'antd';
import type { ActionType } from '@ant-design/pro-table';


const loopTreeItem = (tree: API.MenuList[]): API.MenuList[] =>
    tree.map(({ children, ...item }) => ({
        ...item,
        title: item.name,
        value: item.id,
        children: children && loopTreeItem(children),
    }));


export type PermsFormProps = {
    modalVisible: boolean;
    onChange: (modalVisible: boolean) => void;
    actionRef: React.MutableRefObject<ActionType | undefined>;
    values: API.RoleList | undefined;
};

const PermsForm: React.FC<PermsFormProps> = (props) => {
    const { actionRef, modalVisible, onChange, values } = props;
    const [menuData, setMenuData] = useState<API.MenuList>();
    const [apiData, setApiData] = useState<API.MenuList>();
    const [checkedKeys, setCheckedKeys] = useState([]);
    const [checkedList, setCheckedList] = useState([]);

    const onCheck = (keys: any, info: any) => {
        let allKeys = keys.checked;
        const parentKey = info.node.parent_id;
        if (allKeys.indexOf(parentKey)) {
            setCheckedKeys(allKeys);
        } else {
            allKeys = allKeys.push(parentKey);
            setCheckedKeys(allKeys);
        }
    };

    const onCheckChange = (checkedValues: any) => {
        setCheckedList(checkedValues);
    };

    useEffect(() => {
        queryMenus().then((res) => {
            setMenuData(loopTreeItem(res.data));
        });
        queryApis().then((res) => {
            setApiData(res.data);
        });
        getRolePermsByID(values.id).then((res) => {
            setCheckedKeys(res.data.menus_id);
            setCheckedList(res.data.apis_id);
        });
    }, []);
    return (
        <DrawerForm
            //title="设置权限"
            visible={modalVisible}
            onVisibleChange={onCancel}
            onFinish={() => {
                updatePermsRole(values.id, {
                    menus_id: checkedKeys,
                    apis_id: checkedList,
                })
                    .then((res) => {
                        if (res.code === 200 && res.status === true) {
                            message.success(res.message);
                            if (actionRef.current) {
                                actionRef.current.reload();
                            }
                        }
                    })
                    .then(() => {
                        onCancel();
                    });
                //return true;
            }}
        >
            <h3>菜单权限</h3>
            <Divider />

            <Tree
                checkable
                checkStrictly
                style={{ width: 330 }}
                //defaultCheckedKeys={selectedKeys}
                //defaultSelectedKeys={selectedKeys}
                autoExpandParent={true}
                selectable={false}
                onCheck={onCheck}
                checkedKeys={checkedKeys}
                treeData={menuData}
            />

            <Divider />
            <h3>API权限</h3>
            <Divider />
            <Checkbox.Group style={{ width: '100%' }} value={checkedList} onChange={onCheckChange}>
                {apiData.map((item, index) => {
                    return (
                        <div key={index}>
                            <h4>{item.name}</h4>
                            <Row>
                                {item.children.map((item, index) => {
                                    return (
                                        <Col span={4} key={index}>
                                            <Checkbox value={item.id}>{item.name}</Checkbox>
                                        </Col>
                                    );
                                })}
                            </Row>
                            <Divider />
                        </div>
                    );
                })}
            </Checkbox.Group>
        </DrawerForm>
    );
};

export default PermsForm;