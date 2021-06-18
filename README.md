## 说明
* 十三水的AI出牌算法，速度非常快，一般在微秒级，附带测试方案
* 目前该算法基于golang实现，后续将实现js和lua版本
* 算法基于常归玩法，即52张不带大小王的情况开发

## 优势
* 通常的十三水AI算法都是穷举法，即13张牌遍历出5张，再依次检测是不是特定牌型，速度较慢，因为13张遍历出5张就有（13*12*11*10*9）种组合
* 本算法基于二叉树的思想，将13张牌先拆出5张特定组合来，再将剩下的牌再拆出5张特定的牌型出来
* 具体步骤1：从N张牌中，选出特定的牌型的5张，归到左节点；再将剩余的牌归档到右节点，
* 具体步骤2：将步骤1的右节点递归上述步骤1，这样拆成的3组节点，组成牌型组合
* 具体步骤3：将上述步骤中得到的所有牌型组合进行横向比较，一定是左节点大于右节点的情况，才算正确的牌型组合
* 具体步骤4：将所有正确的牌型组合中，进行纵向比较，选出最优解  
* 本算法速度较快：特殊牌型的计算一般在微秒级，普通牌型因组合较多，计算一般都在10毫秒级，部分较简单的在微秒级  
* 该算法会选出大量组合方案，然后横向比较，选出最优解
* 项目中有测试示例，基本涵盖所有牌型的组合